package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/service"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Start(ctx context.Context, msg *tgbotapi.Message) {
	resp := h.service.Start(ctx, service.StartDto{
		ChatId:    msg.Chat.ID,
		UserId:    msg.From.ID,
		FirstName: msg.From.FirstName,
		ChatType:  msg.Chat.Type,
	})

	b, _ := telegram.BotFromCtx(ctx)

	m := tgbotapi.NewMessage(msg.Chat.ID, resp)
	b.Bot.Send(m)
	return
}
