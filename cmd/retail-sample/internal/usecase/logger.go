package usecase

type logger interface {
	Error(string, string, error)
	Info(string, string)
}
