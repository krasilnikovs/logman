package entity

import "time"

type Server struct {
	Id        int
	Name      string
	Host      string
	LogInfo   LogInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}
