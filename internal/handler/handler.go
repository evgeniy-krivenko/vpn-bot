package handler

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/logger"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/e"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"strings"
	"time"
)

type CommandOrMsgHandlers map[string]func(ctx context.Context, msg *tgbotapi.Message)
type CallbackQueryHandlers map[string]func(ctx context.Context, cq *tgbotapi.CallbackQuery)

const (
	ReqId = "req-id"
)

type Handler struct {
	commandHandlers       CommandOrMsgHandlers
	messageHandlers       CommandOrMsgHandlers
	callbackQueryHandlers CallbackQueryHandlers
	useCases              UseCase
	services              Service
	log                   logger.Logger
}

func NewHandler(uc UseCase, sv Service, l logger.Logger) *Handler {
	return &Handler{useCases: uc, services: sv, log: l}
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
		"terms-confirmed": h.TermsConfirmed,
		"create":          h.CreateConnection,
		"open-connection": h.OpenConnection,
		"activate":        h.ActivateConnection,
		"get-config":      h.GetConfiguration,
		"show-ads":        h.ShowAds,
	}
}

func (h *Handler) HandleCommand(ctx context.Context, msg *tgbotapi.Message) {
	if cb, ok := h.commandHandlers[msg.Command()]; ok {
		go func() {
			defer func() {
				r := recover()
				if r != nil {
					h.log.Errorf(
						"panic with command handler with command: %s",
						msg.Command(),
					)
				}
			}()

			cb(h.setReqId(ctx), msg)
		}()
	}
}

func (h *Handler) HandleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if cb, ok := h.messageHandlers[msg.Text]; ok {
		go func() {
			defer func() {
				r := recover()
				if r != nil {
					h.log.Errorf(
						"panic in handle message goroutine with text: %s",
						msg.Text,
					)
				}
			}()

			cb(h.setReqId(ctx), msg)
		}()
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
		go func() {
			defer func() {
				r := recover()
				if r != nil {
					h.log.Errorf(
						"panic with callbackQueryHandler with pattern: %s",
						pattern,
					)
				}
			}()

			cb(h.setReqId(ctx), cq)
		}()
	}
}

func (h *Handler) sendSelfClearingErrMsg(ctx context.Context, chatId int64, msg string) {
	b, _ := telegram.BotFromCtx(ctx)

	m := tgbotapi.NewMessage(chatId, escape(msg))
	m.ParseMode = Markdown

	respMsg, _ := b.Bot.Send(m)

	select {
	case <-ctx.Done():
		return
	default:
		go func() {
			time.Sleep(15 * time.Second)

			delMsg := tgbotapi.NewDeleteMessage(chatId, respMsg.MessageID)

			_, err := b.Bot.Send(delMsg)
			if err != nil {
				h.log.WithContextReqId(ctx).Errorf(
					"%w",
					e.Warp("error to send delete message", err),
				)
			}
		}()
	}
}

func (h *Handler) setReqId(ctx context.Context) context.Context {
	id := uuid.NewString()
	return context.WithValue(ctx, ReqId, id)
}
