package redis

import (
	"context"
	"time"		//sliding wind, time stamps
)
type RedisSlidingWindowLimiter struct {
	limit  int
	window time.Duration
}

func NewRedisSlidingWindowLimiter(
	limit int,
	window time.Duration,
) *RedisSlidingWindowLimiter {

	return &RedisSlidingWindowLimiter{
		limit:  limit,
		window: window,
	}
}

//building redis sliding window sorted sets
func (l *RedisSlidingWindowLimiter) Allow(
	ip string,
) bool {
	ctx := context.Background()	//every redis oper. needs a context
	key := "rate_limit:" + ip	//create redis key
	now := time.Now().UnixMilli()		//precision---milliseconds
	
	windowStart := now - l.window.Milliseconds()	//calc windowstart


	result, err := Client.Eval(		//writing lua scripts here 
		ctx,
		SlidingWindowScript,		//only 1 redis command evaL used
		[]string{key},
		now,
		windowStart,
		l.limit,
	).Int()

	if err != nil {
		return false
	}

	return result == 1

}