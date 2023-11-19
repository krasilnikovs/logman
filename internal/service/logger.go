package service

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
}
