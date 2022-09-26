package services

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TermsFirstStep  = "terms:1"
	TermsSecondStep = "terms:2"
	TermsThirdStep  = "terms:3"
	TermsForthStep  = "terms:4"
	CommonKeyboards = "CommonKeyboards"
	ErrorKeyboard   = "error"
)

type InlineKeyboards map[string]tgbotapi.InlineKeyboardMarkup

type KeyboardService struct {
}

var inlineMap = InlineKeyboards{
	TermsFirstStep: tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Прочитать", "terms:2"),
		),
	),
	TermsSecondStep: tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<<", "terms:1"),
			tgbotapi.NewInlineKeyboardButtonData("Далее", "terms:3"),
		),
	),
	TermsThirdStep: tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<<", "terms:2"),
			tgbotapi.NewInlineKeyboardButtonData("Далее", "terms:4"),
		),
	),
	TermsForthStep: tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<<", "terms:3"),
			tgbotapi.NewInlineKeyboardButtonData("Согласен ✔️", "terms-confirmed"),
		),
	),
	CommonKeyboards: tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Посмотреть конфигурации", "/configs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список серверов", "/servers"),
		),
	),
	ErrorKeyboard: tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Начать завоно", "/start"),
		),
	),
}

func (s *Service) GetInlineKeyboard(key string) (*tgbotapi.InlineKeyboardMarkup, error) {
	kb, ok := inlineMap[key]
	if !ok {
		return nil, errors.New("wrong inline keyboard key")
	}

	return &kb, nil
}

func (s *Service) NewInlineKeyboard(rows [][]struct{ Text, Data string }) *tgbotapi.InlineKeyboardMarkup {
	if len(rows) == 0 {
		return nil
	}
	var rowsSlice [][]tgbotapi.InlineKeyboardButton

	for _, row := range rows {

		var newRow []tgbotapi.InlineKeyboardButton
		for _, btns := range row {
			btn := tgbotapi.NewInlineKeyboardButtonData(btns.Text, btns.Data)
			newRow = append(newRow, btn)
		}
		rowsSlice = append(rowsSlice, newRow)
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(rowsSlice...)
	return &markup
}
