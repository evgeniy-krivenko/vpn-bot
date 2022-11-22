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
		h.log.WithContextReqId(ctx).
			Errorf("error to get response from use case: %w", err)
		h.sendSelfClearingErrMsg(ctx, msg.Chat.ID, ErrorCommonMessage)
		return
	}

	m := tgbotapi.NewMessage(msg.Chat.ID, escape(response.Msg))
	m.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	keyboard := h.services.NewInlineKeyboard(response.Keys)
	if keyboard != nil {
		m.ReplyMarkup = keyboard
	}

	m.ParseMode = "MarkdownV2"
	b.Bot.Send(m)
	return
}
