// readiness endpoint
package handlers

import (
	"net/http"

	"api-gateway/internal/runtime"
)

func ReadinessHandler(
	manager *runtime.Manager,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		if !manager.Ready(r.Context()) {

			http.Error(
				w,
				"NOT READY",
				http.StatusServiceUnavailable,
			)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	}
}
