package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SplitCallbackCalendarData(callbackData string) (string, string, int, int, int) {
	x := strings.Split(callbackData, "_")
	action := x[1]
	step := x[2]
	year, _ := strconv.Atoi(x[3])
	month, _ := strconv.Atoi(x[4])
	day, _ := strconv.Atoi(x[5])
	return action, step, year, month, day
}

func GetCallbackCalendarData(action string, step string, year int, month int, day int) string {
	return fmt.Sprintf("cbcal_%v_%v_%v_%v_%v", action, step, year, month, day)
}

func BuildYearCalendarWidget(minDate time.Time) tgbotapi.InlineKeyboardMarkup {
	minYear := minDate.Year()
	maxYear := minYear + 3
	currentYear := time.Now().Year()
	showBackNavButton := minYear > currentYear
	var navButtons []tgbotapi.InlineKeyboardButton
	if showBackNavButton {
		navButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<<", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, minYear-1, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(">>", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, maxYear+1, 0, 0)),
		)
	} else {
		navButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("×", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(">>", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, maxYear+1, 0, 0)),
		)
	}

	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v", minYear), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_STEP_MONTH, minYear, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v", minYear+1), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_STEP_MONTH, minYear+1, 0, 0)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v", minYear+2), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_STEP_MONTH, minYear+2, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v", minYear+3), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_STEP_MONTH, minYear+3, 0, 0)),
		),
		navButtons,
	)

	return replyMarkup
}
