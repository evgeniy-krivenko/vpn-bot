package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/e"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"time"
)

func (h *Handler) TermsConfirmed(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	b, _ := telegram.BotFromCtx(ctx)

	resp, err := h.useCases.TermsConfirmed(ctx, int(cq.From.ID))
	if err != nil {
		logrus.Errorf("error confirmed user terms for userId {%d}: %s", cq.From.ID, err.Error())
		h.sendSelfClearingErrMsg(ctx, cq.Message.Chat.ID, ErrorCommonMessage)
		return
	}

	// delete button form prev message
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
			Error(e.Warp("error to send edit message in terms confirmed handler", err))
	}

	newMsg := tgbotapi.NewMessage(cq.Message.Chat.ID, escape(resp.Msg))
	newMsg.ParseMode = "MarkdownV2"
	keyboard := h.services.NewInlineKeyboard(resp.Keys)
	if keyboard != nil {
		newMsg.ReplyMarkup = keyboard
	}

	send = Retry(b.Bot.Send, retryCount, time.Second*retrySecondCount)
	_, err = send(newMsg)
	if err != nil {
		h.log.WithContextReqId(ctx).
			Error(e.Warp("error to send edit message in terms confirmed handler", err))
	}

	h.log.Infof("success confirm terms for user with id %d", cq.From.ID)
}
