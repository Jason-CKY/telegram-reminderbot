package handler

import (
	"github.com/Jason-CKY/telegram-reminderbot/pkg/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func InitializeReminder(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// TODO: Initialize reminder in directus

	// Reply to user message, with keyboard commands to cancel and placeholder text to enter reminder text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, core.REMINDER_BUILDER_MESSAGE)

	msg.ReplyMarkup = core.GetCancelKeyboard()
	msg.ReplyToMessageID = update.Message.MessageID
	// msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}
