package handler

import (
	"github.com/Jason-CKY/telegram-reminderbot/pkg/core"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil { // If we got a message
		if update.Message.IsCommand() {
			HandleCommand(update, bot)
		} else {
			HandleMessage(update, bot)
		}
	}
}

func HandleMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg tgbotapi.MessageConfig
	reminderInConstruction, _ := schemas.GetReminderInConstruction(update.Message.Chat.ID, update.Message.From.ID)

	if update.Message.Text == core.CANCEL_MESSAGE {
		reminder := schemas.Reminder{
			Id:             "",
			ChatId:         update.Message.Chat.ID,
			FromUserId:     update.Message.From.ID,
			FileId:         "",
			Timezone:       "Asia/Singapore",
			Frequency:      "",
			Time:           "",
			ReminderText:   "",
			InConstruction: true,
		}
		err := reminder.DeleteReminderInConstruction()
		if err != nil {
			log.Fatal(err)
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, core.CANCEL_OPERATION_MESSAGE)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	} else if reminderInConstruction != nil {
		// TOOD: https://github.com/Jason-CKY/telegram-reminderbot/blob/main/app/menu.py#L61
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Found reminder text")
		msg.ReplyToMessageID = update.Message.MessageID
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
	}

	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func HandleCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Extract the command from the Message.
	switch update.Message.Command() {
	case "help":
		msg.Text = core.HELP_MESSAGE
	case "start":
		msg.Text = core.HELP_MESSAGE
	case "support":
		msg.Text = core.SUPPORT_MESSAGE
	case "remind":
		InitializeReminder(update, bot)
		return
	case "list":
		msg.Text = "list command handling"
	case "settings":
		msg.Text = "settings command handling"
	default:
		return
	}

	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}
