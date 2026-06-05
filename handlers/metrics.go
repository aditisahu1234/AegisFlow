package handlers

import (
	"encoding/json"
	"net/http"

	"api-gateway/middleware"
)

func MetricsHandler(		//create handler func
	w http.ResponseWriter,
	r *http.Request,
){
	w.Header().Set(		//tell browser it's json
		"Content-Type",
		"application/json",
	)
	json.NewEncoder(w).Encode(		//return metrics
		middleware.GlobalMetrics,
	)
}

//now register route in main.go