package entity

type KeyPath string

type Credential struct {
	Id        int
	Name      string  `validate:"required"`
	Path      KeyPath `validate:"required"`
	CreatedAt string  `validate:"required"`
	UpdatedAt string  `validate:"required"`
}
