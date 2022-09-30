package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/usecases"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartUseCase interface {
	Start(ctx context.Context, dto *entity.User) (*usecases.Response, error)
}

type TermsUseCase interface {
	Terms(ctx context.Context, id int) (*usecases.Response, error)
	TermsConfirmed(ctx context.Context, userId int) (*usecases.Response, error)
}

type ConnectionUseCase interface {
	GetConnections(ctx context.Context, usr *entity.User) (*usecases.ResponseWithKeys, error)
	CreateConnection(ctx context.Context, usr *entity.User, serverId int) ([]usecases.ResponseWithKeys, error)
	// ActivateConnection(ctx context.Context, id int) error
	OpenConnection(ctx context.Context, id int) (*usecases.ResponseWithKeys, error)
}

type UserUseCase interface {
	GetUserById(id int64) (*entity.User, error)
}

type UseCase interface {
	StartUseCase
	TermsUseCase
	ConnectionUseCase
	UserUseCase
}

type KeyboardService interface {
	GetInlineKeyboard(key string) (*tgbotapi.InlineKeyboardMarkup, error)
	NewInlineKeyboard(rows [][]struct{ Text, Data string }) *tgbotapi.InlineKeyboardMarkup
}

type Service interface {
	KeyboardService
}
