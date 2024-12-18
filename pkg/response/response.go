package response

import (
	"encoding/json"
	"net/http"
	"fmt"
)


type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Predefined error responses
var ErrUnauthorized = &ErrorResponse{
    Code:    401,
    Message: "unauthorized",
}
func (e *ErrorResponse) Error() string {
    return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Code:    status,
		Message: message,
	})
}
