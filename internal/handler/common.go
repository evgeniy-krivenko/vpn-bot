package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	ErrorMessageWithStart = "Что-то пошло не так\\. Нажмите /start"
	ErrorCommonMessage    = "Что\\-то пошло не так, попробуйте написать позже"
	Markdown              = "MarkdownV2"
)

func (h *Handler) getIdFromCtx(ctx context.Context) (int, error) {
	payload := ctx.Value("queryPayload").(string)
	id, err := strconv.Atoi(payload)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (h *Handler) checkIsUserExists(id int64) (*entity.User, error) {
	return h.useCases.GetUserById(id)
}

func (h *Handler) responseWithCommonError(ctx context.Context, chatId int64) {
	logrus.Errorf("error to get user for chat %d", chatId)
	b, _ := telegram.BotFromCtx(ctx)
	m := tgbotapi.NewMessage(chatId, ErrorMessageWithStart)
	m.ParseMode = Markdown
	b.Bot.Send(m)
}
