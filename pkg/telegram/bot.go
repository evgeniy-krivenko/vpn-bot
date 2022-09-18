package telegram

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ctxBot struct{}

type Handler interface {
	HandleCommand(ctx context.Context, message *tgbotapi.Message)
	HandleMessage(ctx context.Context, message *tgbotapi.Message)
	HandleCallbackQuery(ctx context.Context, message *tgbotapi.Message, cq *tgbotapi.CallbackQuery)
}

type Bot struct {
	Bot *tgbotapi.BotAPI
	Handler
}

func NewBot(bot *tgbotapi.BotAPI, h Handler) *Bot {
	return &Bot{Bot: bot, Handler: h}
}

func (b *Bot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.Bot.GetUpdatesChan(updateConfig)

	err := b.handleUpdates(updates)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		//if update.Message == nil {
		//	fmt.Println("msg nil")
		//	continue
		//}

		ctx := context.Background()
		ctxWithBot := ContextWithBot(ctx, b)

		if update.CallbackQuery != nil {
			fmt.Println("cb query")
			b.HandleCallbackQuery(ctxWithBot, update.Message, update.CallbackQuery)
			continue
		}

		if update.Message.IsCommand() {
			b.HandleCommand(ctxWithBot, update.Message)
			continue
		}

		b.HandleMessage(ctxWithBot, update.Message)

	}
	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	b.Bot.Send(msg)
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
