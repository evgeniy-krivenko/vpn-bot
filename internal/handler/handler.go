package handler

import (
	"context"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/service"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandlers map[string]func(ctx context.Context, msg *tgbotapi.Message)

type Service interface {
	Start(ctx context.Context, dto service.StartDto) string
}

type Handler struct {
	commandHandlers CommandHandlers
	service         Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitHandlers() {
	ch := CommandHandlers{
		"start": h.Start,
	}
	h.commandHandlers = ch
}

func (h *Handler) HandleCommand(ctx context.Context, msg *tgbotapi.Message) {
	if cb, ok := h.commandHandlers[msg.Command()]; ok {
		go cb(ctx, msg)
	}
}

func (h *Handler) HandleMessage(ctx context.Context, message *tgbotapi.Message) {
	b, err := telegram.BotFromCtx(ctx)
	if err != nil {
		fmt.Printf(err.Error())
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, введите команду")
	b.Bot.Send(msg)
}
