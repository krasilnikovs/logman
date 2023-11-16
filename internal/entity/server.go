package entity

type LogLocation struct {
	Path   string
	Format string
}

type Server struct {
	Id          int
	Name        string
	Host        string
	LogLocation LogLocation
	CreatedAt   string
	UpdatedAt   string
}
