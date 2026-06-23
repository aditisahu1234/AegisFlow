// will own application's readiness state
// redis health separate from application health
package readiness

import "sync/atomic"

var Ready atomic.Bool

func IsReady() bool {
	return Ready.Load()
}

func SetReady(value bool) {
	Ready.Store(value)
}
