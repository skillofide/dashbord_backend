// Package notificationv1 contains the Notification service types.
// The notification service communicates with browsers via WebSocket (not gRPC),
// but this package defines the internal event types exchanged over NATS.
package notificationv1

// Event types published to NATS and forwarded to WebSocket clients.
const (
	EventSubmissionGraded = "SUBMISSION_GRADED"
	EventCodeRunComplete  = "CODE_RUN_COMPLETE"
)

// WebSocketEvent is the payload sent over WebSocket to the browser client.
type WebSocketEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// SubmissionGradedPayload is the payload for EventSubmissionGraded.
type SubmissionGradedPayload struct {
	SubmissionId  string `json:"submission_id"`
	ProblemId     string `json:"problem_id"`
	UserId        string `json:"user_id"`
	OverallStatus string `json:"overall_status"`
	RuntimeMs     int64  `json:"runtime_ms"`
	MemoryKb      int64  `json:"memory_kb"`
	PassedCount   int    `json:"passed_count"`
	TotalCount    int    `json:"total_count"`
}

// CodeRunPayload is the payload for EventCodeRunComplete.
type CodeRunPayload struct {
	JobId        string `json:"job_id"`
	UserId       string `json:"user_id"`
	OverallStatus string `json:"overall_status"`
	CompileError string `json:"compile_error,omitempty"`
}
