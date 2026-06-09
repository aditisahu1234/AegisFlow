package redis

import "sync/atomic"

var RedisHealthy atomic.Bool		//instead of every sinlge req check, keep a shared flag
