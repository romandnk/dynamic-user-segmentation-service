package logger

//go:generate mockgen -source=logger.go -destination=mock/mock.go logger

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
}
