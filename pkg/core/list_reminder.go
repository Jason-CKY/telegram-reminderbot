package core

import (
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

func BuildListReminderMarkup(reminders []schemas.Reminder, page int) (string, tgbotapi.InlineKeyboardMarkup) {
	messageText := ""
	maxDisplayedReminders := page * utils.MAX_REMINDERS_PER_PAGE
	if len(reminders) < maxDisplayedReminders {
		maxDisplayedReminders = len(reminders)
	}
	displayedReminders := reminders[(page-1)*utils.MAX_REMINDERS_PER_PAGE : maxDisplayedReminders]
	var reminderSelectButtons []tgbotapi.InlineKeyboardButton
	for i, reminder := range displayedReminders {
		number := (page-1)*utils.MAX_REMINDERS_PER_PAGE + i + 1
		messageText += fmt.Sprintf(
			"%v%v)    %v (%v at %v)\n",
			utils.REMINDER_PREFIX,
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

	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		reminderSelectButtons,
		navButtons,
	)
	return messageText, replyMarkup
}
