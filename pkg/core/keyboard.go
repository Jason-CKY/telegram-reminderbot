package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetCancelKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var CancelKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🚫 Cancel"),
		),
	)
	CancelKeyboard.InputFieldPlaceholder = "Enter reminder text."
	return CancelKeyboard
}
