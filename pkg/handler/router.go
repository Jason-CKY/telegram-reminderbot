package handler

import (
	"github.com/Jason-CKY/telegram-reminderbot/pkg/core"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	"github.com/google/uuid"
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
	reminderInConstruction, _ := schemas.GetReminderInConstruction(update.Message.Chat.ID, update.Message.From.ID)

	if update.Message.Text == utils.CANCEL_MESSAGE {
		err := reminderInConstruction.DeleteReminderInConstruction()
		if err != nil {
			log.Fatal(err)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, utils.CANCEL_OPERATION_MESSAGE)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	} else if reminderInConstruction != nil {
		// TOOD: https://github.com/Jason-CKY/telegram-reminderbot/blob/main/app/menu.py#L61
		core.BuildReminder(reminderInConstruction, update, bot)
	}
}

func HandleCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Extract the command from the Message.
	switch update.Message.Command() {
	case "help":
		msg.Text = utils.HELP_MESSAGE
	case "start":
		msg.Text = utils.HELP_MESSAGE
	case "support":
		msg.Text = utils.SUPPORT_MESSAGE
	case "remind":
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
		log.Infof("Authorized on account %s", bot.Self.UserName)
		err = reminder.Create()
		if err != nil {
			log.Fatal(err)
		}
		// Reply to user message, with keyboard commands to cancel and placeholder text to enter reminder text
		msg.Text = utils.REMINDER_BUILDER_MESSAGE
		cancelKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(utils.CANCEL_MESSAGE),
			),
		)
		cancelKeyboard.InputFieldPlaceholder = "Enter reminder text."
		msg.ReplyMarkup = cancelKeyboard
		msg.ReplyToMessageID = update.Message.MessageID
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
