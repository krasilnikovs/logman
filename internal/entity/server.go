package entity

const (
	// LogLocationFormatJson is a json format of log location
	LogLocationFormatJson = "json"
)

type LogLocation struct {
	Path   string `validate:"required"`
	Format string `validate:"required,eq=json"`
}

type Server struct {
	Id          int
	Name        string      `validate:"required"`
	Host        string      `validate:"required,hostname|ip"`
	LogLocation LogLocation `validate:"required"`
	CreatedAt   string      `validate:"required"`
	UpdatedAt   string      `validate:"required"`
}
