package usecases

import "github.com/evgeniy-krivenko/particius-vpn-bot/internal/logger"

type UseCase struct {
	*StartUseCase
	*TermsUseCase
	*ConnectionUseCase
	*UserUseCase
}

type UseCaseConfig struct {
	Repository Repository
	Grpc       GrpcService
	Log        logger.Logger
}

func NewUseCase(cnf UseCaseConfig) *UseCase {
	r := cnf.Repository
	grpc := cnf.Grpc
	l := cnf.Log
	return &UseCase{
		StartUseCase:      NewStartUseCase(r),
		TermsUseCase:      NewTermsUseCase(r, grpc),
		ConnectionUseCase: NewConnectionUseCase(r, grpc, l),
		UserUseCase:       NewUserUseCase(r),
	}
}
