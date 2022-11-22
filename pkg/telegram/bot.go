package telegram

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ctxBot struct{}

type Handler interface {
	HandleCommand(ctx context.Context, message *tgbotapi.Message)
	HandleMessage(ctx context.Context, message *tgbotapi.Message)
	HandleCallbackQuery(ctx context.Context, message *tgbotapi.Message, cq *tgbotapi.CallbackQuery)
}

type Bot struct {
	Bot     *tgbotapi.BotAPI
	handler Handler
}

func NewBot(bot *tgbotapi.BotAPI, h Handler) *Bot {
	return &Bot{Bot: bot, handler: h}
}

func (b *Bot) Start(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.Bot.GetUpdatesChan(updateConfig)
	var err error

	go func() {
		err = b.handleUpdates(ctx, updates)
	}()

	go func() {
		<-ctx.Done()
		b.Bot.StopReceivingUpdates()
	}()

	return err
}

func (b *Bot) handleUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ctxWithBot := ContextWithBot(ctx, b)

			if update.CallbackQuery != nil {
				b.handler.HandleCallbackQuery(ctxWithBot, update.Message, update.CallbackQuery)
				continue
			}

			if update.Message.IsCommand() {
				b.handler.HandleCommand(ctxWithBot, update.Message)
				continue
			}

			b.handler.HandleMessage(ctxWithBot, update.Message)
		}
	}
	return nil
}

func ContextWithBot(ctx context.Context, b *Bot) context.Context {
	return context.WithValue(ctx, ctxBot{}, b)
}

func BotFromCtx(ctx context.Context) (*Bot, error) {
	if b, ok := ctx.Value(ctxBot{}).(*Bot); ok {
		return b, nil
	}
	return nil, errors.New("error get bot from context")
}
