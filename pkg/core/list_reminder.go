package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SplitCallbackListReminderData(callbackData string) (string, string, int) {
	x := strings.Split(callbackData, "_")
	action := x[1]
	step := x[2]
	page, _ := strconv.Atoi(x[3])
	return action, step, page
}

func GetCallbackListReminderData(action string, step string, page int) string {
	return fmt.Sprintf("lr_%v_%v_%v", action, step, page)
}

func parseReminderFrequencyToText(reminder schemas.Reminder) string {
	frequencyText := strings.Split(reminder.Frequency, "-")
	frequency := frequencyText[0]
	tz, _ := time.LoadLocation(reminder.Timezone)

	switch frequency {
	case utils.REMINDER_ONCE:
		reminderTime, _ := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%v %v", frequencyText[1], reminder.Time), tz)
		return reminderTime.Format(utils.PRETTY_DATE_FORMAT)
	case utils.REMINDER_DAILY:
		return "every day"
	case utils.REMINDER_WEEKLY:
		weekday, _ := strconv.Atoi(frequencyText[1])
		return fmt.Sprintf("every %v", time.Weekday(weekday))
	case utils.REMINDER_MONTHLY:
		return fmt.Sprintf("%v of every month", frequencyText[1])
	case utils.REMINDER_YEARLY:
		reminderTime, _ := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%v %v", frequencyText[1], reminder.Time), tz)
		return fmt.Sprintf("%v every year", reminderTime.Format(utils.PRETTY_DATE_FORMAT_WITHOUT_YEAR))
	default:
		return ""
	}
}

func BuildListReminderTextAndMarkup(reminders []schemas.Reminder, page int) (string, tgbotapi.InlineKeyboardMarkup, error) {
	messageText := ""
	maxDisplayedReminders := page * utils.MAX_REMINDERS_PER_PAGE
	if len(reminders) < maxDisplayedReminders {
		maxDisplayedReminders = len(reminders)
	}
	displayedReminders := reminders[(page-1)*utils.MAX_REMINDERS_PER_PAGE : maxDisplayedReminders]
	if len(displayedReminders) == 0 {
		return messageText, tgbotapi.InlineKeyboardMarkup{}, errors.New("no reminders in this page")
	}
	var reminderSelectButtons []tgbotapi.InlineKeyboardButton
	for i, reminder := range displayedReminders {
		prefix := utils.REMINDER_PREFIX
		if reminder.FileId != "" {
			prefix = utils.REMINDER_PHOTO_PREFIX
		}
		number := (page-1)*utils.MAX_REMINDERS_PER_PAGE + i + 1
		messageText += fmt.Sprintf(
			"%v%v)    %v (%v at %v)\n",
			prefix,
			number,
			reminder.ReminderText,
			parseReminderFrequencyToText(reminder),
			reminder.Time,
		)
		reminderSelectButtons = append(
			reminderSelectButtons,
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprint(number),
				GetCallbackListReminderData(utils.CALLBACK_SELECT, reminder.Id, page),
			),
		)
	}

	var navButtons []tgbotapi.InlineKeyboardButton
	if page > 1 {
		navButtons = append(
			navButtons,
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("<< Page %v", page-1),
				GetCallbackListReminderData(utils.CALLBACK_GOTO, utils.CALLBACK_NO_ACTION, page-1),
			),
		)
	}
	if len(reminders) > maxDisplayedReminders {
		navButtons = append(
			navButtons,
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Page %v >>", page+1),
				GetCallbackListReminderData(utils.CALLBACK_GOTO, utils.CALLBACK_NO_ACTION, page+1),
			),
		)
	}

	if len(navButtons) > 0 {
		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
			reminderSelectButtons,
			navButtons,
		)
		return messageText, replyMarkup, nil
	}
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		reminderSelectButtons,
	)
	return messageText, replyMarkup, nil

}

func BuildReminderMenuTextAndMarkup(reminder schemas.Reminder) (string, tgbotapi.InlineKeyboardMarkup, error) {
	nextTriggerTime, err := time.ParseInLocation(utils.DIRECTUS_DATETIME_FORMAT, reminder.NextTriggerTime, time.UTC)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	tz, err := time.LoadLocation(reminder.Timezone)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	msgText := fmt.Sprintf(
		"%v\n\n<b>next trigger time:</b>\n%v\n\n<b>Frequency:</b>\n%v at %v",
		reminder.ReminderText,
		nextTriggerTime.In(tz).Format(utils.DATE_AND_TIME_FORMAT),
		parseReminderFrequencyToText(reminder),
		reminder.Time,
	)

	var editButtons []tgbotapi.InlineKeyboardButton
	if reminder.FileId != "" {
		editButtons = append(editButtons,
			tgbotapi.NewInlineKeyboardButtonData(
				"Show Image",
				GetCallbackListReminderData(utils.CALLBACK_SHOW_IMAGE, reminder.Id, 0),
			),
		)
	}

	editButtons = append(editButtons,
		tgbotapi.NewInlineKeyboardButtonData(
			"Delete",
			GetCallbackListReminderData(utils.CALLBACK_DELETE, reminder.Id, 0),
		),
	)

	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		editButtons,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Back to list",
				GetCallbackListReminderData(utils.CALLBACK_GOTO, reminder.Id, 1),
			),
		),
	)

	return msgText, replyMarkup, nil
}
