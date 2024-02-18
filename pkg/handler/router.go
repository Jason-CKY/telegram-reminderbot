package handler

import (
	"github.com/Jason-CKY/telegram-reminderbot/pkg/core"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message.IsCommand() {
		HandleCommand(update, bot)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
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
