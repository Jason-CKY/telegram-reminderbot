package core

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func HandleErrorSendingReminder(reminder schemas.Reminder) error {
	reminderTriggerTime, err := time.ParseInLocation(utils.DIRECTUS_DATETIME_FORMAT, reminder.NextTriggerTime, time.UTC)
	if err != nil {
		return err
	}
	if reminderTriggerTime.Add(24 * time.Hour).Before(time.Now()) {
		err = reminder.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func TriggerReminder(reminder schemas.Reminder, bot *tgbotapi.BotAPI) {
	chatSettings, _, err := schemas.InsertChatSettingsIfNotPresent(reminder.ChatId)
	if err != nil {
		panic(err)
	}

	if reminder.FileId != "" {
		photo_msg := tgbotapi.NewPhoto(
			reminder.ChatId,
			tgbotapi.FileID(reminder.FileId),
		)
		if reminder.ReminderText != "" {
			photo_msg.Caption = fmt.Sprintf("%v%v%v", utils.REMINDER_PREFIX, reminder.ReminderText, utils.RENEW_REMINDER_TEXT)
		} else {
			photo_msg.Caption = fmt.Sprintf("%v%v", utils.REMINDER_PREFIX, utils.RENEW_REMINDER_TEXT)
		}
		photo_msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
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
		if _, err := bot.Request(photo_msg); err != nil {
			log.Error(err)
			err = HandleErrorSendingReminder(reminder)
			if err != nil {
				log.Error(err)
			}
			return
		}
	} else {
		msg := tgbotapi.NewMessage(
			reminder.ChatId,
			fmt.Sprintf("%v%v%v", utils.REMINDER_PREFIX, reminder.ReminderText, utils.RENEW_REMINDER_TEXT),
		)
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
		if _, err := bot.Request(msg); err != nil {
			log.Error(err)
			err = HandleErrorSendingReminder(reminder)
			if err != nil {
				log.Error(err)
			}
			return
		}
	}
	frequencyText := strings.Split(reminder.Frequency, "-")
	frequency := frequencyText[0]
	if frequency == utils.REMINDER_ONCE {
		err := reminder.Delete()
		if err != nil {
			log.Error(err)
			err = HandleErrorSendingReminder(reminder)
			if err != nil {
				log.Error(err)
			}
			return
		}
	} else {
		nextTriggerTime, err := reminder.CalculateNextTriggerTime(chatSettings)
		if err != nil {
			log.Error(err)
			return
		}
		reminder.NextTriggerTime = nextTriggerTime.Format(utils.DIRECTUS_DATETIME_FORMAT)
		err = reminder.Update()
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func ScheduledReminderTrigger(bot *tgbotapi.BotAPI) {
	var wg sync.WaitGroup
	for {
		dueReminders, err := schemas.GetDueReminders()
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(dueReminders); i++ {
			wg.Add(1)
			reminder := dueReminders[i]
			go func(reminder schemas.Reminder, bot *tgbotapi.BotAPI) {
				defer wg.Done()
				TriggerReminder(reminder, bot)
			}(reminder, bot)
		}
		wg.Wait()
		time.Sleep(2 * time.Second)
	}
}
