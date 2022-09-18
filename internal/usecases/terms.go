package usecases

import (
	"context"
	"errors"
	"fmt"
)

type TermsUseCase struct {
	Repository
}

func NewTermsUseCase(r Repository) *TermsUseCase {
	return &TermsUseCase{
		Repository: r,
	}
}

func (c *TermsUseCase) Terms(ctx context.Context, id int) (*Response, error) {
	term, err := c.GetTermsByOrder(id)
	if err != nil {
		return nil, err
	}
	return &Response{
		Msg:         term.Text,
		KeyboardKey: fmt.Sprintf("terms:%d", term.OrderNumber),
	}, nil
}

func (c *TermsUseCase) TermsConfirmed(ctx context.Context, userId int) (*Response, error) {
	user := c.GetUserById(userId)
	if user == nil {
		return nil, errors.New("user is not found")
	}

	user.IsConfirmTerms = true
	err := c.SaveUser(user)
	if err != nil {
		return nil, err
	}

	// TODO - перемесить текст в бд
	return &Response{Msg: "Создать или проверить подключения можно через основное меню    👇", KeyboardKey: ""}, nil
}
