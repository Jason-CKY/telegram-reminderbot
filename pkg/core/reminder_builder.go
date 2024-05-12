package core

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildReminder(reminderInConstruction *schemas.Reminder, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if reminderInConstruction.ReminderText == "" {
		reminderInConstruction.ReminderText = update.Message.Text
		err := reminderInConstruction.Update()
		if err != nil {
			log.Fatal(err)
		}
		msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "enter reminder time in <HH>:<MM> format.")
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	} else if reminderInConstruction.Time == "" {
		time := update.Message.Text
		if !utils.IsValidTime(time) {
			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "Failed to parse time. Please enter time again.")
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		} else {
			reminderInConstruction.Time = time
			err := reminderInConstruction.Update()
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "Once-off reminder or recurring reminder?")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(utils.REMINDER_ONCE),
					tgbotapi.NewKeyboardButton(utils.REMINDER_DAILY),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(utils.REMINDER_WEEKLY),
					tgbotapi.NewKeyboardButton(utils.REMINDER_MONTHLY),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(utils.REMINDER_YEARLY),
					tgbotapi.NewKeyboardButton(utils.CANCEL_MESSAGE),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	} else if reminderInConstruction.Frequency == "" {
		switch update.Message.Text {
		case utils.REMINDER_ONCE:
			reminderInConstruction.Frequency = update.Message.Text
			err := reminderInConstruction.Update()
			if err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "once-off reminder selected.")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}

			// TODO: Monthly Calendar widget
			msg = tgbotapi.NewMessage(reminderInConstruction.ChatId, utils.CALLBACK_CALENDAR_SELECT_YEAR)
			msg.ReplyToMessageID = update.Message.MessageID
			minYear := time.Now().Year()
			msg.ReplyMarkup = BuildYearCalendarWidget(minYear)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}

		case utils.REMINDER_DAILY:
			reminderInConstruction.Frequency = update.Message.Text
			reminderInConstruction.InConstruction = false
			err := reminderInConstruction.Update()
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, fmt.Sprintf("✅ Reminder set for every day at %v", reminderInConstruction.Time))
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		case utils.REMINDER_WEEKLY:
			reminderInConstruction.Frequency = update.Message.Text
			err := reminderInConstruction.Update()
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "Which day of week do you want to set your weekly reminder?")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(utils.MONDAY),
					tgbotapi.NewKeyboardButton(utils.TUESDAY),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(utils.WEDNESDAY),
					tgbotapi.NewKeyboardButton(utils.THURSDAY),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(utils.FRIDAY),
					tgbotapi.NewKeyboardButton(utils.CANCEL_MESSAGE),
				),
			)
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		case utils.REMINDER_MONTHLY:
			reminderInConstruction.Frequency = update.Message.Text
			err := reminderInConstruction.Update()
			if err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "Which day of the month do you want to set your monthly reminder? (1-31)")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		case utils.REMINDER_YEARLY:
			reminderInConstruction.Frequency = update.Message.Text
			err := reminderInConstruction.Update()
			if err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "yearly reminder selected.")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
			// TODO: Calendar widget
		}
	} else if reminderInConstruction.Frequency == utils.REMINDER_WEEKLY {
		val, ok := utils.DAY_OF_WEEK[update.Message.Text]
		if ok {
			reminderInConstruction.Frequency = fmt.Sprintf("%v-%v", utils.REMINDER_WEEKLY, val)
			reminderInConstruction.InConstruction = false
			err := reminderInConstruction.Update()
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, fmt.Sprintf("✅ Reminder set for every %v at %v", update.Message.Text, reminderInConstruction.Time))
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	} else if reminderInConstruction.Frequency == utils.REMINDER_MONTHLY {
		// day_of_month := update.Message.Text
		day_of_month, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			return
		}
		if day_of_month >= 1 && day_of_month <= 31 {
			reminderInConstruction.Frequency = fmt.Sprintf("%v-%v", utils.REMINDER_WEEKLY, day_of_month)
			reminderInConstruction.InConstruction = false
			err := reminderInConstruction.Update()
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, fmt.Sprintf("✅ Reminder set for every %v of every month at %v", utils.ParseDayOfMonth(day_of_month), reminderInConstruction.Time))
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		} else {
			msg := tgbotapi.NewMessage(reminderInConstruction.ChatId, "Invalid day of month [1-31]")
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}

	}
}
