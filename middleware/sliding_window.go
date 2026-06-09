package middleware

import (
	"net"
	"net/http"
	"sync"		//mutex to avoid race condition
	"time"
	"fmt"
)

type Limiter interface {		//create an interface
	Allow(ip string) bool
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


//allow function, which allows or blocks
func (l *SlidingWindowLimiter) Allow(ip string) bool {


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
		GlobalMetrics.BlockedRequests++
		fmt.Println("BLOCKED:", ip)
		return false				//req blocked
	}

	validRequests = append(validRequests, now)		//add current req, if not blocked
	l.requests[ip] = validRequests		//save back to map

	fmt.Println("ALLOWED:", ip)
	GlobalMetrics.AllowedRequests++
	return true //req accepted
}

//middleware signature
func SlidingWindowMiddleware(

	limiter Limiter,		//sliding window limiter object, limiter stuff from main.go gets passed here
	next http.Handler,					// returns http.Handler
) http.Handler {

	return http.HandlerFunc(func(			//creates a handler func
		w http.ResponseWriter,				//every req passes through this code
		r *http.Request,
	) {

		//extract client ip address
		ip, _, err := net.SplitHostPort(
			r.RemoteAddr,
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

		if !limiter.Allow(ip) {		//ask limiter, true or false??
			http.Error(
				w,
				"rate limit exceeded",		//block request
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)	//limiter blcoked req, call actual endpoint
	})
}
