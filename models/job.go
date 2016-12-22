package models

const (
	SCHEDULED = "scheduled"
	FINISHED  = "finished"
	ERROR     = "error"
)

type Job struct {
	ID         string `json:"id,"`
	Status     string `json:"status"`
	Executable string `json:"executable"`
	Schedule   string `json:"schedule"`
	Type       string `json:"type"`
}

// storage for all jobs for now. This should be stored in a persistent datastore. But for sake of prototype, it will stored here.
var Jobs []Job

// Counter to generate Job ID. This can be change by generating a UUID
var JobCounter = 0
