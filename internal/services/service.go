package services

type Service struct {
	*KeyboardService
	*Crypto
}

func New() *Service {
	return &Service{
		KeyboardService: new(KeyboardService),
		Crypto:          new(Crypto),
	}
}
