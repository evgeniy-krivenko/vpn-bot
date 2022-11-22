package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/e"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	retryCount        = 3
	retrySecondCount  = 3
	userWaitAdsSecond = 10
)

func (h *Handler) GetConnections(ctx context.Context, msg *tgbotapi.Message) {
	user, err := h.checkIsUserExists(msg.From.ID)
	if err != nil {
		h.log.WithContextReqId(ctx).
			Errorf("error get conns, user %d is not exists in database", msg.From.ID)
		h.sendSelfClearingErrMsg(ctx, msg.Chat.ID, ErrorMessageWithStart)
		return
	}

	b, _ := telegram.BotFromCtx(ctx)

	response, err := h.useCases.GetConnections(ctx, user)
	if err != nil {
		h.log.WithContextReqId(ctx).
			Errorf("error to get connections for user %d", msg.From.ID)
		h.sendSelfClearingErrMsg(ctx, msg.Chat.ID, ErrorCommonMessage)
		return
	}

	m := tgbotapi.NewMessage(msg.Chat.ID, escape(response.Msg))
	m.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	keyboard := h.services.NewInlineKeyboard(response.Keys)
	if keyboard != nil {
		m.ReplyMarkup = keyboard
	}

	m.ParseMode = Markdown
	send := Retry(b.Bot.Send, 3, time.Second*3)
	if _, err := send(m); err != nil {
		h.log.WithContextReqId(ctx).
			Error(e.Warp("error to send message", err))
	}
	return
}

func (h *Handler) CreateConnection(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	user, err := h.checkIsUserExists(cq.From.ID)
	if err != nil {
		h.log.WithContextReqId(ctx).
			Errorf("error create conn, user %d is not exists in database", cq.From.ID)
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorMessageWithStart)
		return
	}

	id, err := h.getIdFromCtx(ctx)
	if err != nil {
		h.log.WithContextReqId(ctx).
			Error(e.Warp("error parse id from ctx when create conn", err))
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorCommonMessage)
		return
	}

	b, _ := telegram.BotFromCtx(ctx)

	response, err := h.useCases.CreateConnection(ctx, user, id)
	if err != nil {
		h.log.Errorf("error to create connections for user %d: %s", cq.From.ID, err.Error())
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorCommonMessage)
		return
	}

	for _, r := range response {
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, escape(r.Msg))
		keyboard := h.services.NewInlineKeyboard(r.Keys)
		if keyboard != nil {
			m.ReplyMarkup = keyboard
		}

		m.ParseMode = Markdown

		b.Bot.Send(m)
	}
	return
}

func (h *Handler) OpenConnection(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	id, err := h.getIdFromCtx(ctx)
	b, _ := telegram.BotFromCtx(ctx)

	if err != nil {
		h.log.Errorf("error parse id from ctx for chat id %d: %s", cq.Message.Chat.ID, err.Error())
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorMessageWithStart)
		return
	}

	resp, err := h.useCases.OpenConnection(ctx, id)
	if err != nil {
		h.log.Errorf("error to open connection with id %d: %s", id, err.Error())
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorCommonMessage)
		return
	}

	m := tgbotapi.NewMessage(cq.Message.Chat.ID, escape(resp.Msg))

	keyboard := h.services.NewInlineKeyboard(resp.Keys)
	if keyboard != nil {
		m.ReplyMarkup = keyboard
	}

	m.ParseMode = Markdown
	if msg, err := b.Bot.Send(m); err != nil {
		h.log.Errorf("error to send message with id %d: %w", msg.MessageID, err)
	}
}

// ActivateConnection - call after show ads for user
func (h *Handler) ActivateConnection(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	id, err := h.getIdFromCtx(ctx)
	b, _ := telegram.BotFromCtx(ctx)

	if err != nil {
		h.log.Errorf("error parse id from ctx for chat id %d: %w", cq.Message.Chat.ID, err)
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorMessageWithStart)
		return
	}

	resp, err := h.useCases.ActivateConnection(ctx, id)
	if err != nil {
		h.log.Errorf("error to activate connection with id %d: %s", id, err.Error())
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorCommonMessage)
		return
	}

	// delete button from prev message
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				h.log.WithContextReqId(ctx).
					Error("error in show ads waiting goroutine\n", r)
			}
		}()

		m := tgbotapi.NewEditMessageText(
			cq.Message.Chat.ID,
			cq.Message.MessageID,
			escape(cq.Message.Text),
		)
		m.ParseMode = Markdown

		send := Retry(b.Bot.Send, retryCount, time.Second*retrySecondCount)
		_, err = send(m)
		if err != nil {
			h.log.WithContextReqId(ctx).
				Error(e.Warp("error to send edit message in activate connections handler", err))
		}
	}()

	m := tgbotapi.NewMessage(cq.Message.Chat.ID, escape(resp.Msg))

	keyboard := h.services.NewInlineKeyboard(resp.Keys)
	if keyboard != nil {
		m.ReplyMarkup = keyboard
	}

	m.ParseMode = Markdown
	if msg, err := b.Bot.Send(m); err != nil {
		h.log.Errorf("error to send message with id %d: %s", msg.MessageID, err.Error())
	}
	return
}

func (h *Handler) GetConfiguration(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	id, err := h.getIdFromCtx(ctx)

	b, _ := telegram.BotFromCtx(ctx)

	if err != nil {
		h.log.Errorf(
			"error parse id from ctx for chat id %d: %w",
			cq.Message.Chat.ID,
			err,
		)
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorMessageWithStart)
		return
	}

	response, err := h.useCases.GetConfiguration(ctx, id)

	if err != nil {
		h.log.Errorf(
			"error to get configuration for user %d with id %d: %w",
			cq.From.ID,
			id,
			err,
		)
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorCommonMessage)
		return
	}

	for _, r := range response {
		m := tgbotapi.NewMessage(cq.Message.Chat.ID, escape(r.Msg))
		keyboard := h.services.NewInlineKeyboard(r.Keys)
		if keyboard != nil {
			m.ReplyMarkup = keyboard
		}

		m.ParseMode = Markdown

		respMsg, _ := b.Bot.Send(m)

		if r.IsMessageDelete {
			go func() {
				select {
				case <-ctx.Done():
					return
				case <-time.After(30 * time.Second):
					delMsg := tgbotapi.NewDeleteMessage(cq.Message.Chat.ID, respMsg.MessageID)

					b.Bot.Send(delMsg)
				}
			}()
		}
	}
	return
}

// ShowAds - for showing user advertising before activate connection
func (h *Handler) ShowAds(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	user, err := h.checkIsUserExists(cq.From.ID)
	if err != nil {
		h.log.WithContextReqId(ctx).
			Errorf("error create conn, user %d is not exists in database", cq.From.ID)
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorMessageWithStart)
		return
	}

	id, err := h.getIdFromCtx(ctx)
	if err != nil {
		h.log.WithContextReqId(ctx).
			Error(e.Warp("error parse id from ctx when create conn", err))
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorCommonMessage)
		return
	}

	b, _ := telegram.BotFromCtx(ctx)

	r, _ := h.useCases.ShowAds(ctx, user, id)

	m := tgbotapi.NewMessage(cq.Message.Chat.ID, escape(r.Msg))
	var respMsg tgbotapi.Message
	m.ParseMode = Markdown
	send := Retry(b.Bot.Send, retryCount, time.Second*retrySecondCount)
	if respMsg, err = send(m); err != nil {
		h.log.WithContextReqId(ctx).
			Errorf("[ShowAds] error to send message: %w", err)
		return
	}

	go func() {
		defer func() {
			r := recover()
			if r != nil {
				h.log.WithContextReqId(ctx).
					Error("error in show ads waiting goroutine\n", r)
			}
		}()
		keyboard := h.services.NewInlineKeyboard(r.Keys)
		time.Sleep(userWaitAdsSecond * time.Second)
		m := tgbotapi.NewEditMessageTextAndMarkup(
			cq.Message.Chat.ID,
			respMsg.MessageID,
			escape(r.Msg),
			*keyboard,
		)
		m.ParseMode = Markdown

		send := Retry(b.Bot.Send, retryCount, time.Second*retrySecondCount)
		if respMsg, err = send(m); err != nil {
			h.log.WithContextReqId(ctx).
				Error(e.Warp("error to send edit message in show ads handler", err))
		}
	}()

}
