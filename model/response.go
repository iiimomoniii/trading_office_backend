package model

import "time"

// BaseResponse — standard response wrapper
type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ErrorDetail — error detail payload
type ErrorDetail struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	ErrorCode string    `json:"errorCode"`
	Error     string    `json:"error"`
	MessageEN string    `json:"messageEN"`
	MessageTH string    `json:"messageTH"`
	Path      string    `json:"path"`
}

// PaginatedResponse — for list endpoints
type PaginatedResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Total   int64  `json:"total"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
}
