package services

type Service struct {
	*KeyboardService
}

func New() *Service {
	return &Service{new(KeyboardService)}
}
