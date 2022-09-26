package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/usecases"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UseCase interface {
	Start(ctx context.Context, dto *entity.User) (*usecases.Response, error)
	Terms(ctx context.Context, id int) (*usecases.Response, error)
	TermsConfirmed(ctx context.Context, userId int) (*usecases.Response, error)
	GetConnections(ctx context.Context, usr *entity.User) (*usecases.ResponseWithKeys, error)
	CreateConnection(ctx context.Context, usr *entity.User, serverId int) ([]usecases.ResponseWithKeys, error)
	GetUserById(id int64) (*entity.User, error)
}

type KeyboardService interface {
	GetInlineKeyboard(key string) (*tgbotapi.InlineKeyboardMarkup, error)
	NewInlineKeyboard(rows [][]struct{ Text, Data string }) *tgbotapi.InlineKeyboardMarkup
}

type Service interface {
	KeyboardService
}
