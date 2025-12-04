package shared

import (
	"time"
)

type CodeRequest struct {
	ID          string    `json:"id"`
	Language    string    `json:"language"`
	Code        string    `json:"code"`
	Timestamp   time.Time `json:"timestamp"`
}

type CodeResponse struct {
    ID         string    `json:"id"`
    Output     string    `json:"output"`
    Error      string    `json:"error,omitempty"`
    ExitCode   int       `json:"exit_code"`
    Duration   int64     `json:"duration"` // milliseconds
    MemoryUsed int64     `json:"memory_used,omitempty"`
    Timestamp  time.Time `json:"timestamp"`
}