package entity

import "time"

// Supported log formats
const (
	JsonFormat = "json"
)

// LogInfo describes knowledge about log
type LogInfo struct {
	Id   int
	Name string
	// the Location field contains path where logs store
	Location string
	// The Format field contains which logs format used
	Format    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
