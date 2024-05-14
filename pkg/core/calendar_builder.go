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

func BuildDayCalendarWidget(callbackData string) tgbotapi.InlineKeyboardMarkup {
	// TODO: fit the days button on the correct day of the week
	_, _, selectedYear, selectedMonth, _ := SplitCallbackCalendarData(callbackData)
	selectedDate := time.Date(selectedYear, time.Month(selectedMonth), 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Now()

	var dayButtons []tgbotapi.InlineKeyboardButton
	var buttons [][]tgbotapi.InlineKeyboardButton

	legendButtons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("M", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_DAY, selectedYear, selectedMonth, 0)),
		tgbotapi.NewInlineKeyboardButtonData("T", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_DAY, selectedYear, selectedMonth, 0)),
		tgbotapi.NewInlineKeyboardButtonData("W", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_DAY, selectedYear, selectedMonth, 0)),
		tgbotapi.NewInlineKeyboardButtonData("T", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_DAY, selectedYear, selectedMonth, 0)),
		tgbotapi.NewInlineKeyboardButtonData("F", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_DAY, selectedYear, selectedMonth, 0)),
		tgbotapi.NewInlineKeyboardButtonData("S", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_DAY, selectedYear, selectedMonth, 0)),
		tgbotapi.NewInlineKeyboardButtonData("S", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_DAY, selectedYear, selectedMonth, 0)),
	)
	buttons = append(buttons, legendButtons)

	daysInMonth := utils.DaysInMonth(selectedDate)
	for i := 1; i <= utils.CALENDAR_DAY_NUM_ROWS*utils.CALENDAR_DAY_NUM_COLS; i++ {
		if (i < currentDate.Day() && selectedDate.Compare(currentDate) <= 0) || i > daysInMonth {
			dayButtons = append(dayButtons, tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_MONTH, 0, 0, 0)))
		} else {
			dayButtons = append(dayButtons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(i), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_SELECT_DAY, selectedYear, i, 0)))
		}
	}

	for row := 0; row < utils.CALENDAR_DAY_NUM_ROWS; row++ {
		var rowButtons []tgbotapi.InlineKeyboardButton
		for col := 0; col < utils.CALENDAR_DAY_NUM_COLS; col++ {
			rowButtons = append(rowButtons, dayButtons[row*utils.CALENDAR_DAY_NUM_COLS+col])

		}
		buttons = append(buttons, rowButtons)
	}

	showBackNavButton := selectedDate.Compare(currentDate) > 0
	var navButtons []tgbotapi.InlineKeyboardButton
	prevMonth := selectedMonth - 1
	prevYear := selectedYear
	nextMonth := selectedMonth + 1
	nextYear := selectedYear
	if prevMonth < 1 {
		prevMonth = 12
		prevYear = selectedYear - 1
	}
	if nextMonth > 12 {
		nextMonth = 1
		nextYear = selectedYear + 1
	}
	if showBackNavButton {
		navButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<<", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_DAY, prevYear, prevMonth, 0)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v %v", time.Month(selectedMonth).String()[:3], selectedDate.Year()), GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_MONTH, selectedYear, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(">>", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_DAY, nextYear, nextMonth, 0)),
		)
	} else {
		navButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("×", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_DAY, prevYear, selectedMonth, 0)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v %v", time.Month(selectedMonth).String()[:3], selectedDate.Year()), GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_MONTH, selectedYear, 0, 0)),
			tgbotapi.NewInlineKeyboardButtonData(">>", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_DAY, nextYear, nextMonth, 0)),
		)
	}

	buttons = append(buttons, navButtons)
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)

	return replyMarkup
}
