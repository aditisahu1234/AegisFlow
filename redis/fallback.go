package redis

import (
	"time"

	"api-gateway/middleware"
)

//uses the original simple sliding wondow limiter
var FallbackLimiter =
	middleware.NewSlidingWindowLimiter(
		50,		//fallback limit only 50 requests/minute
		time.Minute,
	)