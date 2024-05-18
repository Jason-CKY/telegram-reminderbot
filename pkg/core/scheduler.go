package core

import (
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
			msg := tgbotapi.NewMessage(reminder.ChatId, reminder.ReminderText)
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
