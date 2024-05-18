package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func ScheduledReminderTrigger(bot *tgbotapi.BotAPI) {
	for {
		dueReminders, err := schemas.GetDueReminders()
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(dueReminders); i++ {
			reminder := dueReminders[i]
			msg := tgbotapi.NewMessage(reminder.ChatId, fmt.Sprintf("%v%v%v", utils.REMINDER_PREFIX, reminder.ReminderText, utils.RENEW_REMINDER_TEXT))
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("15m", utils.RENEW_REMINDER_15M),
					tgbotapi.NewInlineKeyboardButtonData("30m", utils.RENEW_REMINDER_30M),
					tgbotapi.NewInlineKeyboardButtonData("1h", utils.RENEW_REMINDER_1H),
					tgbotapi.NewInlineKeyboardButtonData("3h", utils.RENEW_REMINDER_3H),
					tgbotapi.NewInlineKeyboardButtonData("1d", utils.RENEW_REMINDER_1D),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Enter Time", utils.RENEW_REMINDER_CUSTOM),
					tgbotapi.NewInlineKeyboardButtonData("Cancel", utils.RENEW_REMINDER_CANCEL),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
			frequencyText := strings.Split(reminder.Frequency, "-")
			frequency := frequencyText[0]
			if frequency == utils.REMINDER_ONCE {
				err := reminder.DeleteById()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				nextTriggerTime, err := reminder.CalculateNextTriggerTime()
				if err != nil {
					log.Fatal(err)
				}
				reminder.NextTriggerTime = nextTriggerTime.Format(utils.DIRECTUS_DATETIME_FORMAT)
				err = reminder.Update()
				if err != nil {
					log.Fatal(err)
				}
			}

		}
		time.Sleep(2 * time.Second)
	}
}
