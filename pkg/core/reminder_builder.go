package core

import (
	"log"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildReminder(reminderInConstruction *schemas.Reminder, update *tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "test")
	if reminderInConstruction.ReminderText == "" {
		reminderInConstruction.ReminderText = update.Message.Text
		err := reminderInConstruction.Update()
		if err != nil {
			log.Fatal(err)
		}
		msg.Text = "enter reminder time in <HH>:<MM> format."
	}
	msg.ReplyToMessageID = update.Message.MessageID
	return msg
}
