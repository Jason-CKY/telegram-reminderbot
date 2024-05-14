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

func BuildYearCalendarWidget(callbackData string) tgbotapi.InlineKeyboardMarkup {
	_, _, minYear, _, _ := SplitCallbackCalendarData(callbackData)
	maxYear := minYear + 3
	currentYear := time.Now().Year()

	var yearButtons []tgbotapi.InlineKeyboardButton
	var buttons [][]tgbotapi.InlineKeyboardButton
	for year := minYear; year <= maxYear; year++ {
		if year < currentYear {
			yearButtons = append(yearButtons, tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)))
		} else {
			yearButtons = append(yearButtons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v", year), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_STEP_MONTH, year, 0, 0)))
		}

	}

	for row := 0; row < utils.CALENDAR_YEAR_NUM_ROWS; row++ {
		var rowButtons []tgbotapi.InlineKeyboardButton
		for col := 0; col < utils.CALENDAR_YEAR_NUM_COLS; col++ {
			rowButtons = append(rowButtons, yearButtons[row*utils.CALENDAR_YEAR_NUM_COLS+col])
		}
		buttons = append(buttons, rowButtons)
	}

	showBackNavButton := minYear > currentYear
	var navButtons []tgbotapi.InlineKeyboardButton
	if showBackNavButton {
		navButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<<", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, minYear-4, 0, 0)),
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

	buttons = append(buttons, navButtons)
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)

	return replyMarkup
}

func BuildMonthCalendarWidget(callbackData string) tgbotapi.InlineKeyboardMarkup {
	_, _, selectedYear, _, _ := SplitCallbackCalendarData(callbackData)
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	var monthButtons []tgbotapi.InlineKeyboardButton
	var buttons [][]tgbotapi.InlineKeyboardButton

	for i := 1; i <= 12; i++ {
		if i < int(currentMonth) && selectedYear <= currentYear {
			monthButtons = append(monthButtons, tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_MONTH, 0, 0, 0)))
		} else {
			monthButtons = append(monthButtons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v", time.Month(i).String()[:3]), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_SELECT_DAY, selectedYear, i, 0)))
		}
	}

	for row := 0; row < utils.CALENDAR_MONTH_NUM_ROWS; row++ {
		var rowButtons []tgbotapi.InlineKeyboardButton
		for col := 0; col < utils.CALENDAR_MONTH_NUM_COLS; col++ {
			rowButtons = append(rowButtons, monthButtons[row*utils.CALENDAR_MONTH_NUM_COLS+col])
		}
		buttons = append(buttons, rowButtons)
	}

	showBackNavButton := selectedYear > currentYear
	var navButtons []tgbotapi.InlineKeyboardButton
	if showBackNavButton {
		navButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<<", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_MONTH, selectedYear-1, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(selectedYear), GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, selectedYear, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(">>", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_MONTH, selectedYear+1, 0, 0)),
		)
	} else {
		navButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("×", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_MONTH, selectedYear-1, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(currentYear), GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, selectedYear, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(">>", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_MONTH, selectedYear+1, 0, 0)),
		)
	}

	buttons = append(buttons, navButtons)
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)

	return replyMarkup
}
