package middleware

import (
	"net"
	"net/http"
	"sync"		//mutex to avoid race condition
	"time"
	"fmt"
	"context"
	"api-gateway/telemetry"
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Limiter interface {		//create an interface
	Allow(
		ctx context.Context,
		ip string,
	) bool

}

// create limiter struct

type SlidingWindowLimiter struct {
	mu       sync.Mutex				//locks access to shared data
	window   time.Duration
	limit    int
	requests map[string][]time.Time		//gives the timestamp, for ip addresses
}


//constructor
func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:    limit,
		window:   window,
		requests: make(map[string][]time.Time),
	}
}
//method for active ip count
func (l *SlidingWindowLimiter) ActiveIPCount() int {

    l.mu.Lock()
    defer l.mu.Unlock()

    return len(l.requests)
}

//allow function, which allows or blocks
func (l *SlidingWindowLimiter) Allow(
	ctx context.Context,
	ip string,
	) bool {


	fmt.Println("Request received from:", ip)
    //bool gives true-->allow or false-->block request
	
	l.mu.Lock()		//lock the data
	defer l.mu.Unlock()

	GlobalMetrics.TotalRequests++		//updating global metrics for total number of requests

	now := time.Now()		//get current time
	timestamps := l.requests[ip]		//check ip's request history

	var validRequests []time.Time		//fresh list create, remove old timestamps

	for _, t := range timestamps {			//remove old requests, only requests from last 1 minute remain, (set accordingly)
		if now.Sub(t) < l.window {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= l.limit {	//check if limit of requests exceeded
		l.requests[ip] = validRequests
		telemetry.ActiveIPs.Store(
			int64(
				len(l.requests),
			),
		)
		GlobalMetrics.BlockedRequests++
		fmt.Println("BLOCKED:", ip)
		return false				//req blocked
	}

	validRequests = append(validRequests, now)		//add current req, if not blocked
	l.requests[ip] = validRequests		//save back to map
	telemetry.ActiveIPs.Store(
		int64(
			len(l.requests),
		),
	)
	fmt.Println("ALLOWED:", ip)
	GlobalMetrics.AllowedRequests++
	return true //req accepted
}



/*-------------------------------------------------------------------------------------
----------------------------------------------------------------------------------------
				ENTRY POINT FOR EVERY HTTP REQUEST
------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------
*/



//middleware signature 
func SlidingWindowMiddleware(

	limiter Limiter,		//sliding window limiter object, limiter stuff from main.go gets passed here
	next http.Handler,					// returns http.Handler
) http.Handler {

	return http.HandlerFunc(func(			//creates a handler func
		w http.ResponseWriter,				//every req passes through this code
		r *http.Request,
	) {

		ctx, span :=		//ROOT SPAN
			telemetry.Tracer.Start(
				r.Context(),
				"http_request",
			)

		defer span.End()

		r = r.WithContext(ctx)

		span.SetAttributes(		//adding SPAN attributes
			attribute.String(
				"http.method",
				r.Method,
			),
		
			attribute.String(
				"http.path",
				r.URL.Path,
			),
		)

		start := time.Now()

		telemetry.		//at request entry
			ActiveRequestsCounter.
			Add(
				r.Context(),
				1,
			)
		defer func() {

			telemetry.
				ActiveRequestsCounter.
				Add(
					r.Context(),
					-1,
				)
			
			duration := time.Since(start)
			
			telemetry.
				RequestDurationHistogram.
				Record(
					r.Context(),
					duration.Seconds(),
				)
		}()
		telemetry.RequestCounter.Add(		//updating telemetry metrics
			context.Background(),
			1,
		)

		//extract client ip address
		ip, _, err := net.SplitHostPort(
			r.RemoteAddr,
		)

		span.SetAttributes(
			attribute.String(
				"client.ip",
				ip,
			),
		)
						//all requests count towards the same bucket
		if err != nil {
			http.Error(
				w,
				"invalid client address",
				http.StatusInternalServerError,
			)
			return
		}

		span.AddEvent(		//add events
			"rate_limit_check_started",
		)

		if !limiter.Allow(
			r.Context(),		//active trace context travels downstream
			ip,
		) {		//ask limiter, true or false??

			span.AddEvent(
				"rate_limit_exceeded",
			)

			span.RecordError(	//error recording
				errors.New(
					"rate limit exceeded",
				),
			)
			
			span.SetStatus(
				codes.Error,
				"rate limit exceeded",
			)

			http.Error(
				w,
				"rate limit exceeded",		//block request
				http.StatusTooManyRequests,
			)
			return
		}
		span.AddEvent(
			"rate_limit_check_passed",
		)

		next.ServeHTTP(w, r)	//limiter blcoked req, call actual endpoint
		
	})
}
