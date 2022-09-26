package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type CommandOrMsgHandlers map[string]func(ctx context.Context, msg *tgbotapi.Message)
type CallbackQueryHandlers map[string]func(ctx context.Context, cq *tgbotapi.CallbackQuery)

type Handler struct {
	commandHandlers       CommandOrMsgHandlers
	messageHandlers       CommandOrMsgHandlers
	callbackQueryHandlers CallbackQueryHandlers
	useCases              UseCase
	services              Service
}

func NewHandler(uc UseCase, sv Service) *Handler {
	return &Handler{useCases: uc, services: sv}
}

func (h *Handler) HandleCommand(ctx context.Context, msg *tgbotapi.Message) {
	if cb, ok := h.commandHandlers[msg.Command()]; ok {
		go cb(ctx, msg)
	}
}

func (h *Handler) HandleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if cb, ok := h.messageHandlers[msg.Text]; ok {
		go cb(ctx, msg)
	}
}

func (h *Handler) HandleCallbackQuery(ctx context.Context, msg *tgbotapi.Message, cq *tgbotapi.CallbackQuery) {
	var pattern string
	splitData := strings.Split(cq.Data, ":")

	if len(splitData) > 1 {
		pattern = splitData[0]
		ctx = context.WithValue(ctx, "queryPayload", splitData[1])
	} else {
		pattern = cq.Data
	}

	if cb, ok := h.callbackQueryHandlers[pattern]; ok {
		go cb(ctx, cq)
	}
}

// InitHandlers define handlers here
func (h *Handler) InitHandlers() {
	h.commandHandlers = CommandOrMsgHandlers{
		"start":       h.Start,
		"connections": h.GetConnections,
	}
	h.messageHandlers = CommandOrMsgHandlers{
		getConfig:    h.GetConfig,
		checkConfigs: h.CheckConfigs,
	}

	h.callbackQueryHandlers = CallbackQueryHandlers{
		"terms":           h.Terms,
		"terms-confirmed": h.TermsConfirmed,
		"create":          h.CreateConnections,
	}
}
