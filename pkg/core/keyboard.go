package core

import (
	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetCancelKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var CancelKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(utils.CANCEL_MESSAGE),
		),
	)
	CancelKeyboard.InputFieldPlaceholder = "Enter reminder text."
	return CancelKeyboard
}
