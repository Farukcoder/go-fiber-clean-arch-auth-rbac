package domain

import "time"

// RequestLog is the DB entity for HTTP request audit logs.
type RequestLog struct {
	ID           int64     `json:"id"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	StatusCode   int       `json:"status_code"`
	DurationMs   float64   `json:"duration_ms"`
	RequestBody  string    `json:"request_body"`
	ResponseBody string    `json:"response_body"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	ErrorType    string    `json:"error_type"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}
