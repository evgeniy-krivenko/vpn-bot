package usecases

type UseCase struct {
	*StartUseCase
	*TermsUseCase
	*ConnectionUseCase
}

type UseCaseConfig struct {
	Repository Repository
	Service    Service
}

func NewUseCase(cnf UseCaseConfig) *UseCase {
	r := cnf.Repository
	s := cnf.Service
	return &UseCase{
		StartUseCase:      NewStartUseCase(r),
		TermsUseCase:      NewTermsUseCase(r),
		ConnectionUseCase: NewConnectionUseCase(r, s),
	}
}
