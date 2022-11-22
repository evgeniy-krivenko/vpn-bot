package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
)

type Response struct {
	Msg         string
	KeyboardKey string
}

type StartUseCase struct {
	userRepo UserRepository
}

func NewStartUseCase(r Repository) *StartUseCase {
	return &StartUseCase{userRepo: r}
}

func (uc *StartUseCase) Start(ctx context.Context, dto *entity.User) (*ResponseWithKeys, error) {
	user := uc.userRepo.GetUserById(int(dto.UserId))
	if user == nil {
		return uc.handleNewUser(ctx, dto)
	}

	return uc.handleExistUser(ctx, user)
}

func (uc *StartUseCase) GetUserById(id int64) (*entity.User, error) {
	user := uc.userRepo.GetUserById(int(id))

	if user == nil {
		return nil, errors.New(fmt.Sprintf("user with id %d is not exists", id))
	}

	return user, nil
}

func (uc *StartUseCase) handleNewUser(ctx context.Context, dto *entity.User) (*ResponseWithKeys, error) {
	if _, err := uc.userRepo.CreateNewUser(dto); err != nil {
		return nil, err
	}

	return uc.getTerms(), nil
}

func (uc *StartUseCase) handleExistUser(ctx context.Context, user *entity.User) (*ResponseWithKeys, error) {
	if user.IsConfirmTerms {
		return &ResponseWithKeys{Msg: MainMenuText}, nil
	}

	return uc.getTerms(), nil
}

func (uc *StartUseCase) getTerms() *ResponseWithKeys {
	resp := ResponseWithKeys{Msg: StartText}

	resp.AddRow(
		resp.AddButton("Подтвердить", "terms-confirmed"),
	)

	return &resp
}
