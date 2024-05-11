package handler

import (
	"github.com/Jason-CKY/telegram-reminderbot/pkg/core"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func InitializeReminder(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	reminder := schemas.Reminder{
		Id:             uuid.New().String(),
		ChatId:         update.Message.Chat.ID,
		FromUserId:     update.Message.From.ID,
		FileId:         "",
		Timezone:       "Asia/Singapore",
		Frequency:      "",
		Time:           "",
		ReminderText:   "",
		InConstruction: true,
	}
	// delete previous reminders in construction to create a new one
	err := reminder.DeleteReminderInConstruction()
	if err != nil {
		log.Fatal(err)
	}
	// create a new reminder
	err = reminder.Create()
	if err != nil {
		log.Fatal(err)
	}
	// Reply to user message, with keyboard commands to cancel and placeholder text to enter reminder text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, core.REMINDER_BUILDER_MESSAGE)

	msg.ReplyMarkup = core.GetCancelKeyboard()
	msg.ReplyToMessageID = update.Message.MessageID
	// msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}
