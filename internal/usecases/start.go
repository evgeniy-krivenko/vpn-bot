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
	Repository
}

func NewStartUseCase(r Repository) *StartUseCase {
	return &StartUseCase{Repository: r}
}

// Start -
func (uc *StartUseCase) Start(ctx context.Context, dto *entity.User) (*Response, error) {
	user := uc.Repository.GetUserById(int(dto.UserId))
	if user == nil {
		return uc.handleNewUser(ctx, dto)
	}

	return uc.handleExistUser(ctx, user)
}

func (uc *StartUseCase) GetUserById(id int64) (*entity.User, error) {
	user := uc.Repository.GetUserById(int(id))

	if user == nil {
		return nil, errors.New(fmt.Sprintf("user with id %d is not exists", id))
	}

	return user, nil
}

func (uc *StartUseCase) handleNewUser(ctx context.Context, dto *entity.User) (*Response, error) {
	if _, err := uc.Repository.CreateNewUser(dto); err != nil {
		return nil, err
	}

	return uc.getTerms()
}

func (uc *StartUseCase) handleExistUser(ctx context.Context, user *entity.User) (*Response, error) {
	fmt.Println(user)
	if user.IsConfirmTerms {
		return &Response{Msg: MainMenuText, KeyboardKey: ""}, nil
	}

	return uc.getTerms()
}

func (uc *StartUseCase) getTerms() (*Response, error) {
	term, err := uc.Repository.GetTermsByOrder(1)
	if err != nil {
		return nil, err
	}

	return &Response{
		Msg:         term.Text,
		KeyboardKey: fmt.Sprintf("terms:%d", term.OrderNumber),
	}, nil
}
