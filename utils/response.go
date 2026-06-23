package utils

import (
	"time"
)

type MetaData struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Meta    MetaData    `json:"meta"`
}

func BuildSuccess(code, message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
		Meta: MetaData{
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}
}

func BuildError(code, message string, errors interface{}) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
		Errors:  errors,
		Meta: MetaData{
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}
}
