package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetConnections(ctx context.Context, msg *tgbotapi.Message) {
	user, err := h.checkIsUserExists(msg.From.ID)
	if err != nil {
		h.responseWithCommonError(ctx, msg.Chat.ID)
		return
	}

	b, _ := telegram.BotFromCtx(ctx)

	response, err := h.useCases.GetConnections(ctx, user)

	if err != nil {
		logrus.Errorf("error to get connections for user %d", msg.From.ID)
		resp := ErrorCommonMessage
		m := tgbotapi.NewMessage(msg.Chat.ID, resp)
		m.ParseMode = Markdown
		b.Bot.Send(m)
		return
	}

	m := tgbotapi.NewMessage(msg.Chat.ID, response.Msg)
	m.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	keyboard := h.services.NewInlineKeyboard(response.Keys)
	if keyboard != nil {
		m.ReplyMarkup = keyboard
	}

	m.ParseMode = Markdown
	b.Bot.Send(m)
	return
}

func (h *Handler) CreateConnections(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	user, err := h.checkIsUserExists(cq.From.ID)
	if err != nil {
		h.responseWithCommonError(ctx, cq.Message.Chat.ID)
		return
	}

	id, err := h.getIdFromCtx(ctx)

	b, _ := telegram.BotFromCtx(ctx)

	if err != nil {
		logrus.Errorf("error parse id from ctx for chat id %d: %s", cq.Message.Chat.ID, err.Error())
		resp := ErrorCommonMessage
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, resp)
		m.ParseMode = Markdown
		b.Bot.Send(m)
		return
	}

	response, err := h.useCases.CreateConnection(ctx, user, id)

	if err != nil {
		logrus.Errorf("error to create connections for user %d: %s", cq.From.ID, err.Error())
		resp := ErrorCommonMessage
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, resp)
		m.ParseMode = Markdown
		b.Bot.Send(m)
		return
	}

	for i, r := range response {
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, r.Msg)
		keyboard := h.services.NewInlineKeyboard(r.Keys)
		if keyboard != nil {
			m.ReplyMarkup = keyboard
		}
		if i != 0 {
			// pass for first message
			m.ParseMode = Markdown
		}

		b.Bot.Send(m)
	}
}

func (h *Handler) OpenConnection(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	id, err := h.getIdFromCtx(ctx)
	b, _ := telegram.BotFromCtx(ctx)

	if err != nil {
		logrus.Errorf("error parse id from ctx for chat id %d: %s", cq.Message.Chat.ID, err.Error())
		resp := ErrorCommonMessage
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, resp)
		m.ParseMode = Markdown
		b.Bot.Send(m)
		return
	}

	resp, err := h.useCases.OpenConnection(ctx, id)
	if err != nil {
		logrus.Errorf("error to open connection with id %d: %s", id, err.Error())
		resp := ErrorCommonMessage
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, resp)
		m.ParseMode = Markdown
		b.Bot.Send(m)
		return
	}

	m := tgbotapi.NewMessage(cq.Message.Chat.ID, resp.Msg)

	keyboard := h.services.NewInlineKeyboard(resp.Keys)
	if keyboard != nil {
		m.ReplyMarkup = keyboard
	}

	m.ParseMode = Markdown
	if msg, err := b.Bot.Send(m); err != nil {
		logrus.Errorf("error to send message with id %d: %s", msg.MessageID, err.Error())
	}
	return
}

func (h *Handler) ActivateConnection(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	id, err := h.getIdFromCtx(ctx)
	b, _ := telegram.BotFromCtx(ctx)

	if err != nil {
		logrus.Errorf("error parse id from ctx for chat id %d: %s", cq.Message.Chat.ID, err.Error())
		resp := ErrorCommonMessage
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, resp)
		m.ParseMode = Markdown
		b.Bot.Send(m)
		return
	}

	resp, err := h.useCases.ActivateConnection(ctx, id)
	if err != nil {
		logrus.Errorf("error to activate connection with id %d: %s", id, err.Error())
		resp := ErrorCommonMessage
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, resp)
		m.ParseMode = Markdown
		b.Bot.Send(m)
		return
	}

	m := tgbotapi.NewMessage(cq.Message.Chat.ID, resp.Msg)

	keyboard := h.services.NewInlineKeyboard(resp.Keys)
	if keyboard != nil {
		m.ReplyMarkup = keyboard
	}

	m.ParseMode = Markdown
	if msg, err := b.Bot.Send(m); err != nil {
		logrus.Errorf("error to send message with id %d: %s", msg.MessageID, err.Error())
	}
	return
}
