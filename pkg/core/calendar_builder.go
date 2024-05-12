package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
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

func BuildYearCalendarWidget(minYear int) tgbotapi.InlineKeyboardMarkup {
	maxYear := minYear + 3
	currentYear := time.Now().Year()
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

	var yearButtons []tgbotapi.InlineKeyboardButton
	numColumns := 2
	numRows := 2
	var buttons [][]tgbotapi.InlineKeyboardButton
	for year := minYear; year <= maxYear; year++ {
		log.Info(year)
		if year < currentYear {
			yearButtons = append(yearButtons, tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)))
		} else {
			yearButtons = append(yearButtons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v", year), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_STEP_MONTH, year, 0, 0)))
		}

	}

	for row := 0; row < numRows; row++ {
		var rowButtons []tgbotapi.InlineKeyboardButton
		for col := 0; col < numColumns; col++ {
			rowButtons = append(rowButtons, yearButtons[row+col])
		}
		buttons = append(buttons, rowButtons)
	}

	buttons = append(buttons, navButtons)
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)

	return replyMarkup
}

// func BuildMonthCalendarWidget() tgbotapi.InlineKeyboardMarkup {
// 	maxYear := minYear + 3
// 	currentYear := time.Now().Year()
// 	showBackNavButton := minYear > currentYear
// 	var navButtons []tgbotapi.InlineKeyboardButton
// 	if showBackNavButton {
// 		navButtons = tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("<<", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, minYear-4, 0, 0)),
// 			tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)),
// 			tgbotapi.NewInlineKeyboardButtonData(">>", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, maxYear+1, 0, 0)),
// 		)
// 	} else {
// 		navButtons = tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("×", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)),
// 			tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)),
// 			tgbotapi.NewInlineKeyboardButtonData(">>", GetCallbackCalendarData(utils.CALLBACK_GOTO, utils.CALLBACK_CALENDAR_STEP_YEAR, maxYear+1, 0, 0)),
// 		)
// 	}

// 	var yearButtons []tgbotapi.InlineKeyboardButton
// 	numColumns := 2
// 	numRows := 2
// 	var buttons [][]tgbotapi.InlineKeyboardButton
// 	for year := minYear; year <= maxYear; year++ {
// 		log.Info(year)
// 		if year < currentYear {
// 			yearButtons = append(yearButtons, tgbotapi.NewInlineKeyboardButtonData(" ", GetCallbackCalendarData(utils.CALLBACK_NO_ACTION, utils.CALLBACK_CALENDAR_STEP_YEAR, 0, 0, 0)))
// 		} else {
// 			yearButtons = append(yearButtons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v", year), GetCallbackCalendarData(utils.CALLBACK_SELECT, utils.CALLBACK_CALENDAR_STEP_MONTH, year, 0, 0)))
// 		}

// 	}

// 	for row := 0; row < numRows; row++ {
// 		var rowButtons []tgbotapi.InlineKeyboardButton
// 		for col := 0; col < numColumns; col++ {
// 			rowButtons = append(rowButtons, yearButtons[row+col])
// 		}
// 		buttons = append(buttons, rowButtons)
// 	}

// 	buttons = append(buttons, navButtons)
// 	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
// 		buttons...,
// 	)

// 	return replyMarkup
// }
