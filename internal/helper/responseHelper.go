package helper

import (
	"bookingService/internal/structs"
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(structs.ErrorResponse{
		Error: errorMessage,
	})
}
