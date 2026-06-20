package handlers
import (
    "encoding/json"
    "net/http"		//since http.ResponseWriterr and http.Request live here
)

func HealthHandler(			//handler func
    w http.ResponseWriter,	
    r *http.Request,
){
	w.Header().Set(		//to tell browser response type
		"Content-Type",
		"application/json",
	)
	
	w.WriteHeader(http.StatusOK)
	response := map[string]string{		//create response
		"status": "healthy",
	}
	
	json.NewEncoder(w).Encode(response)	//conver to json nd send to browser for display
}

//register route of handler in main.go