package usecases

type UseCase struct {
	*StartUseCase
	*TermsUseCase
}

func NewUseCase(r Repository) *UseCase {
	return &UseCase{
		StartUseCase: NewStartUseCase(r),
		TermsUseCase: NewTermsUseCase(r),
	}
}
