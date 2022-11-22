package handler

import (
	"context"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

const (
	ErrorMessageWithStart = "Что-то пошло не так. Нажмите /start"
	ErrorCommonMessage    = "Что-то пошло не так, мы уже в курсе и работаем над решением проблемы, попробуйте позже"
	Markdown              = "MarkdownV2"
)

func (h *Handler) getIdFromCtx(ctx context.Context) (int, error) {
	payload := ctx.Value("queryPayload").(string)
	id, err := strconv.Atoi(payload)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (h *Handler) getReqIdFromCtx(ctx context.Context) string {
	reqId, ok := ctx.Value(ReqId).(string)
	if !ok {
		reqId = "none"
	}

	return reqId
}

func (h *Handler) checkIsUserExists(id int64) (*entity.User, error) {
	return h.useCases.GetUserById(id)
}

type SendFunc func(c tgbotapi.Chattable) (tgbotapi.Message, error)

func Retry(f SendFunc, retries int, delay time.Duration) SendFunc {
	return func(c tgbotapi.Chattable) (tgbotapi.Message, error) {
		for r := 0; ; r++ {
			resp, err := f(c)
			if err == nil || r >= retries {
				return resp, err
			}

			select {
			case <-time.After(delay):
				return tgbotapi.Message{}, err
			}
		}
	}
}

func escape(msg string) string {
	var escSymbols = []string{"_", "*", "[", "]", "(", ")", "~", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, s := range escSymbols {
		msg = strings.ReplaceAll(msg, s, fmt.Sprintf("\\%s", s))
	}
	return msg
}
