package usecases

type UserUseCase struct {
	repo Repository
}

func NewUserUseCase(repo Repository) *UserUseCase {
	return &UserUseCase{repo}
}
