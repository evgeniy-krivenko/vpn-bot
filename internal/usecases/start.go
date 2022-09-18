package usecases

import (
	"context"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/services"
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

// Start получить данные о юзере, которые мы сохраним в базе
func (uc *StartUseCase) Start(ctx context.Context, dto entity.User) (*Response, error) {
	user := uc.Repository.GetUserById(int(dto.UserId))
	if user == nil {
		return uc.handleNewUser(ctx, dto)
	}

	return uc.handleExistUser(ctx, dto)
}

func (uc *StartUseCase) handleNewUser(ctx context.Context, dto entity.User) (*Response, error) {
	if _, err := uc.Repository.CreateNewUser(&dto); err != nil {
		return nil, err
	}
	// TODO: надо продумать пайплайн сообщений из пользовательских соглашений
	// возможно нужно будет храть в базе таблицу с прочтением того или иного сообщения
	term, err := uc.Repository.GetTermsByOrder(1)
	if err != nil {
		return nil, err
	}
	return &Response{
		Msg:         term.Text,
		KeyboardKey: fmt.Sprintf("terms:%d", term.OrderNumber),
	}, nil
}

func (uc *StartUseCase) handleExistUser(ctx context.Context, dto entity.User) (*Response, error) {
	return &Response{Msg: "Что хотите сделать?", KeyboardKey: services.CommonKeyboards}, nil
}
