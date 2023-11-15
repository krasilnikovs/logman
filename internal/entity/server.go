package entity

import "time"

type LogLocation struct {
	Path   string
	Format string
}

type Server struct {
	Id          int
	Name        string
	Host        string
	LogLocation LogLocation
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
