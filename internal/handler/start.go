package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Start(ctx context.Context, msg *tgbotapi.Message) {
	b, _ := telegram.BotFromCtx(ctx)

	response, err := h.useCases.Start(ctx, &entity.User{
		ChatId:    msg.Chat.ID,
		UserId:    msg.From.ID,
		FirstName: msg.From.FirstName,
		ChatType:  msg.Chat.Type,
	})

	if err != nil {
		// logging
		resp := "Что\\-то пошло не так, попробуйте написать позже"
		m := tgbotapi.NewMessage(msg.Chat.ID, resp)
		b.Bot.Send(m)
		return
	}

	m := tgbotapi.NewMessage(msg.Chat.ID, response.Msg)
	m.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	keyboard, _ := h.services.GetInlineKeyboard(response.KeyboardKey)
	if keyboard != nil {
		m.ReplyMarkup = keyboard
	}

	m.ParseMode = "MarkdownV2"
	b.Bot.Send(m)
	return
}
