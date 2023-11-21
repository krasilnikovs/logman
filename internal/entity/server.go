package entity

const (
	// LogLocationFormatJson is a json format of log location
	LogLocationFormatJson = "json"
)

type LogFolderPath string
type LogFormat string

type Server struct {
	Id            int
	Name          string        `validate:"required"`
	Host          string        `validate:"required,hostname|ip"`
	LogFolderPath LogFolderPath `validate:"required"`
	LogFormat     LogFormat     `validate:"required,eq=json"`
	CredentialId  int           `validate:"required"`
	CreatedAt     string        `validate:"required"`
	UpdatedAt     string        `validate:"required"`
}
