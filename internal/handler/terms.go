package handler

import (
	"context"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"strconv"
)

func (h *Handler) Terms(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	b, _ := telegram.BotFromCtx(ctx)
	var newMsg tgbotapi.EditMessageTextConfig

	fmt.Println("terms handler")

	payload := ctx.Value("queryPayload").(string)
	id, err := strconv.Atoi(payload)
	if err != nil {
		logrus.Errorf("error convert payload: %s", err.Error())
		b.Bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, "Что-то пошло не так. Нажмите /start"))
		return
	}

	resp, err := h.useCases.Terms(ctx, id)
	if err != nil {
		logrus.Errorf("error get term from usecase: %s", err.Error())
		b.Bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, "Что-то пошло не так. Нажмите /start"))
		return
	}

	kb, err := h.services.GetInlineKeyboard(resp.KeyboardKey)
	if err != nil {
		logrus.Errorf("error get keyboard: %s", err.Error())
		b.Bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, "Что-то пошло не так. Нажмите /start"))
		return
	}

	newMsg = tgbotapi.NewEditMessageTextAndMarkup(cq.Message.Chat.ID, cq.Message.MessageID, resp.Msg, kb)
	newMsg.ParseMode = "MarkdownV2"
	b.Bot.Send(newMsg)
}

func (h *Handler) TermsConfirmed(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	b, _ := telegram.BotFromCtx(ctx)

	resp, err := h.useCases.TermsConfirmed(ctx, int(cq.From.ID))
	if err != nil {
		logrus.Errorf("error confirmed user terms for userId {%d}: %s", cq.From.ID, err.Error())
		b.Bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, "Что-то пошло не так. Нажмите /start"))
		return
	}

	delMsg := tgbotapi.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)

	b.Bot.Send(delMsg)

	newMsg := tgbotapi.NewMessage(cq.Message.Chat.ID, resp.Msg)
	b.Bot.Send(newMsg)
	newMsg.ParseMode = "MarkdownV2"
	logrus.Infof("success confirm terms for user with id {%d}", cq.From.ID)
}
