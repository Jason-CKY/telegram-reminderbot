package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/core"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil { // If we got a message
		chatSettings, err := schemas.InsertChatSettingsIfNotPresent(update.Message.Chat.ID)
		if err != nil {
			log.Fatal(err)
		}
		if update.Message.IsCommand() {
			HandleCommand(update, bot, chatSettings)
		} else {
			HandleMessage(update, bot, chatSettings)
		}
	} else if update.CallbackQuery != nil {
		chatSettings, err := schemas.InsertChatSettingsIfNotPresent(update.CallbackQuery.Message.Chat.ID)
		if err != nil {
			log.Fatal(err)
		}
		HandleCallbackQuery(update, bot, chatSettings)
	}
}

func HandleMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI, chatSettings *schemas.ChatSettings) {
	reminderInConstruction, _ := schemas.GetReminderInConstruction(update.Message.Chat.ID, update.Message.From.ID)

	if update.Message.Text == utils.CANCEL_MESSAGE {
		if reminderInConstruction != nil {
			err := reminderInConstruction.DeleteReminderInConstruction()
			if err != nil {
				log.Fatal(err)
			}
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, utils.CANCEL_OPERATION_MESSAGE)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	} else if update.Message.Text == utils.SETTINGS_CHANGE_TIMEZONE && chatSettings.Updating {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, utils.CHANGE_TIMEZONE_MESSAGE)
		cancelKeyboard := tgbotapi.NewOneTimeReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(utils.CANCEL_MESSAGE),
			),
		)
		cancelKeyboard.InputFieldPlaceholder = "Enter timezone"
		msg.ReplyMarkup = cancelKeyboard
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}

	} else if reminderInConstruction != nil {
		core.BuildReminder(reminderInConstruction, update, bot)
	}
}

func HandleCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI, chatSettings *schemas.ChatSettings) {

	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Extract the command from the Message.
	switch update.Message.Command() {
	case "help":
		msg.Text = utils.HELP_MESSAGE
	case "start":
		msg.Text = utils.HELP_MESSAGE
	case "support":
		msg.Text = utils.SUPPORT_MESSAGE
	case "remind":
		reminder := schemas.Reminder{
			Id:              uuid.New().String(),
			ChatId:          update.Message.Chat.ID,
			FromUserId:      update.Message.From.ID,
			FileId:          "",
			Timezone:        chatSettings.Timezone,
			Frequency:       "",
			Time:            "",
			ReminderText:    "",
			InConstruction:  true,
			NextTriggerTime: "",
		}
		// delete previous reminders in construction to create a new one
		err := reminder.DeleteReminderInConstruction()
		if err != nil {
			log.Fatal(err)
		}
		// create a new reminder
		err = reminder.Create()
		if err != nil {
			log.Fatal(err)
		}
		// Reply to user message, with keyboard commands to cancel and placeholder text to enter reminder text
		msg.Text = utils.REMINDER_BUILDER_MESSAGE
		cancelKeyboard := tgbotapi.NewOneTimeReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(utils.CANCEL_MESSAGE),
			),
		)
		cancelKeyboard.InputFieldPlaceholder = "Enter reminder text."
		msg.ReplyMarkup = cancelKeyboard
		msg.ReplyToMessageID = update.Message.MessageID
	case "list":
		// TODO
		msg.Text = "list command handling"
	case "settings":
		tz, _ := time.LoadLocation(chatSettings.Timezone)
		msg.Text = fmt.Sprintf("<b>Your current settings:</b>\n\n- timezone: %v\n- local time: %v", chatSettings.Timezone, time.Now().In(tz).Format(utils.DATE_AND_TIME_FORMAT_WITHOUT_YEAR))
		msg.ParseMode = "html"
		msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(utils.SETTINGS_CHANGE_TIMEZONE),
				tgbotapi.NewKeyboardButton(utils.CANCEL_MESSAGE),
			),
		)
	default:
		return
	}

	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func HandleCallbackQuery(update *tgbotapi.Update, bot *tgbotapi.BotAPI, chatSettings *schemas.ChatSettings) {
	reminderInConstruction, _ := schemas.GetReminderInConstruction(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.ID)
	if strings.HasPrefix(update.CallbackQuery.Data, "cbcal") && reminderInConstruction != nil {
		action, step, _, _, _ := core.SplitCallbackCalendarData(update.CallbackQuery.Data)
		tz, _ := time.LoadLocation(reminderInConstruction.Timezone)
		if action != utils.CALLBACK_NO_ACTION {
			if action == utils.CALLBACK_GOTO {
				if step == utils.CALLBACK_CALENDAR_STEP_YEAR {
					// user clicks on navigation button on year view
					replyMarkup := core.BuildYearCalendarWidget(update.CallbackQuery.Data, tz)
					editedMessage := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID,
						utils.CALLBACK_CALENDAR_SELECT_YEAR,
						replyMarkup,
					)
					if _, err := bot.Request(editedMessage); err != nil {
						log.Fatal(err)
					}
				} else if step == utils.CALLBACK_CALENDAR_STEP_MONTH {
					// user clicks on navigation button on month view
					replyMarkup := core.BuildMonthCalendarWidget(update.CallbackQuery.Data, tz)
					editedMessage := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID,
						utils.CALLBACK_CALENDAR_SELECT_MONTH,
						replyMarkup,
					)
					if _, err := bot.Request(editedMessage); err != nil {
						log.Fatal(err)
					}
				} else if step == utils.CALLBACK_CALENDAR_STEP_DAY {
					// user clicks on navigation button on day view
					replyMarkup := core.BuildDayCalendarWidget(update.CallbackQuery.Data, tz)
					editedMessage := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID,
						utils.CALLBACK_CALENDAR_SELECT_DAY,
						replyMarkup,
					)
					if _, err := bot.Request(editedMessage); err != nil {
						log.Fatal(err)
					}
				}
			} else if action == utils.CALLBACK_SELECT {
				if step == utils.CALLBACK_CALENDAR_STEP_YEAR {
					// user clicks on a year
					replyMarkup := core.BuildMonthCalendarWidget(update.CallbackQuery.Data, tz)
					editedMessage := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID,
						utils.CALLBACK_CALENDAR_SELECT_MONTH,
						replyMarkup,
					)
					if _, err := bot.Request(editedMessage); err != nil {
						log.Fatal(err)
					}
				} else if step == utils.CALLBACK_CALENDAR_STEP_MONTH {
					// user clicks on a month
					replyMarkup := core.BuildDayCalendarWidget(update.CallbackQuery.Data, tz)
					editedMessage := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID,
						utils.CALLBACK_CALENDAR_SELECT_DAY,
						replyMarkup,
					)
					if _, err := bot.Request(editedMessage); err != nil {
						log.Fatal(err)
					}
				} else if step == utils.CALLBACK_CALENDAR_STEP_DAY {
					// user clicks on a day
					if reminderInConstruction.Frequency == utils.REMINDER_ONCE || reminderInConstruction.Frequency == utils.REMINDER_YEARLY {
						_, _, selectedYear, selectedMonth, selectedDay := core.SplitCallbackCalendarData(update.CallbackQuery.Data)
						// reminderTime stored in db is in UTC, while the date string is in user's timezone, so we need to correct that
						reminderHour, reminderMinute := utils.ParseReminderTime(reminderInConstruction.Time)
						reminderDate := time.Date(selectedYear, time.Month(selectedMonth), selectedDay, reminderHour, reminderMinute, 0, 0, tz)

						replyMessageText := fmt.Sprintf("✅ Reminder set for %v", reminderDate.Format(utils.DATE_AND_TIME_FORMAT))
						if reminderInConstruction.Frequency == utils.REMINDER_ONCE {
							reminderInConstruction.Frequency = fmt.Sprintf("%v-%v", utils.REMINDER_ONCE, reminderDate.Format(utils.DATE_FORMAT))
						} else if reminderInConstruction.Frequency == utils.REMINDER_YEARLY {
							reminderInConstruction.Frequency = fmt.Sprintf("%v-%v", utils.REMINDER_YEARLY, reminderDate.Format(utils.DATE_FORMAT))
							replyMessageText = fmt.Sprintf("✅ Reminder set for every year at %v", reminderDate.Format(utils.DATE_AND_TIME_FORMAT_WITHOUT_YEAR))
						}

						nextTriggerTime, err := reminderInConstruction.CalculateNextTriggerTime()
						if err != nil {
							log.Fatal(err)
						}
						reminderInConstruction.NextTriggerTime = nextTriggerTime.Format(utils.DIRECTUS_DATETIME_FORMAT)
						reminderInConstruction.InConstruction = false
						err = reminderInConstruction.Update()
						if err != nil {
							log.Fatal(err)
						}
						editedMessage := tgbotapi.NewEditMessageText(
							update.CallbackQuery.Message.Chat.ID,
							update.CallbackQuery.Message.MessageID,
							replyMessageText,
						)
						if _, err := bot.Request(editedMessage); err != nil {
							log.Fatal(err)
						}

					}
				}
			}
		}
	} else if strings.HasPrefix(update.CallbackQuery.Data, "renew") && reminderInConstruction == nil {
		reminderText := update.CallbackQuery.Message.Text[:len(update.CallbackQuery.Message.Text)-len(utils.RENEW_REMINDER_TEXT)]
		tz, err := time.LoadLocation(chatSettings.Timezone)
		if err != nil {
			log.Fatal(err)
		}
		switch update.CallbackQuery.Data {
		case utils.RENEW_REMINDER_15M:
			nextTriggerTime := time.Now().In(tz).Add(15 * time.Minute)
			reminder := schemas.Reminder{
				Id:              uuid.New().String(),
				ChatId:          update.CallbackQuery.Message.Chat.ID,
				FromUserId:      update.CallbackQuery.From.ID,
				FileId:          "",
				Timezone:        chatSettings.Timezone,
				Frequency:       fmt.Sprintf("%v-%v", utils.REMINDER_ONCE, nextTriggerTime.Format(utils.DATE_FORMAT)),
				Time:            nextTriggerTime.Format(utils.TIME_ONLY_FORMAT),
				ReminderText:    strings.TrimPrefix(reminderText, utils.REMINDER_PREFIX),
				InConstruction:  false,
				NextTriggerTime: nextTriggerTime.In(time.UTC).Format(utils.DIRECTUS_DATETIME_FORMAT),
			}
			err = reminder.Create()
			if err != nil {
				log.Fatal(err)
			}
			editedMessage := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				fmt.Sprintf("%v\n\nI will remind you again on %v", reminderText, nextTriggerTime.Format(utils.DATE_AND_TIME_FORMAT)),
			)
			if _, err := bot.Request(editedMessage); err != nil {
				log.Fatal(err)
			}
		case utils.RENEW_REMINDER_30M:
			nextTriggerTime := time.Now().In(tz).Add(30 * time.Minute)
			reminder := schemas.Reminder{
				Id:              uuid.New().String(),
				ChatId:          update.CallbackQuery.Message.Chat.ID,
				FromUserId:      update.CallbackQuery.From.ID,
				FileId:          "",
				Timezone:        chatSettings.Timezone,
				Frequency:       fmt.Sprintf("%v-%v", utils.REMINDER_ONCE, nextTriggerTime.Format(utils.DATE_FORMAT)),
				Time:            nextTriggerTime.Format(utils.TIME_ONLY_FORMAT),
				ReminderText:    strings.TrimPrefix(reminderText, utils.REMINDER_PREFIX),
				InConstruction:  false,
				NextTriggerTime: nextTriggerTime.In(time.UTC).Format(utils.DIRECTUS_DATETIME_FORMAT),
			}
			err = reminder.Create()
			if err != nil {
				log.Fatal(err)
			}
			editedMessage := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				fmt.Sprintf("%v\n\nI will remind you again on %v", reminderText, nextTriggerTime.Format(utils.DATE_AND_TIME_FORMAT)),
			)
			if _, err := bot.Request(editedMessage); err != nil {
				log.Fatal(err)
			}
		case utils.RENEW_REMINDER_1H:
			nextTriggerTime := time.Now().In(tz).Add(1 * time.Hour)
			reminder := schemas.Reminder{
				Id:              uuid.New().String(),
				ChatId:          update.CallbackQuery.Message.Chat.ID,
				FromUserId:      update.CallbackQuery.From.ID,
				FileId:          "",
				Timezone:        chatSettings.Timezone,
				Frequency:       fmt.Sprintf("%v-%v", utils.REMINDER_ONCE, nextTriggerTime.Format(utils.DATE_FORMAT)),
				Time:            nextTriggerTime.Format(utils.TIME_ONLY_FORMAT),
				ReminderText:    strings.TrimPrefix(reminderText, utils.REMINDER_PREFIX),
				InConstruction:  false,
				NextTriggerTime: nextTriggerTime.In(time.UTC).Format(utils.DIRECTUS_DATETIME_FORMAT),
			}
			err = reminder.Create()
			if err != nil {
				log.Fatal(err)
			}
			editedMessage := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				fmt.Sprintf("%v\n\nI will remind you again on %v", reminderText, nextTriggerTime.Format(utils.DATE_AND_TIME_FORMAT)),
			)
			if _, err := bot.Request(editedMessage); err != nil {
				log.Fatal(err)
			}
		case utils.RENEW_REMINDER_3H:
			nextTriggerTime := time.Now().In(tz).Add(3 * time.Hour)
			reminder := schemas.Reminder{
				Id:              uuid.New().String(),
				ChatId:          update.CallbackQuery.Message.Chat.ID,
				FromUserId:      update.CallbackQuery.From.ID,
				FileId:          "",
				Timezone:        chatSettings.Timezone,
				Frequency:       fmt.Sprintf("%v-%v", utils.REMINDER_ONCE, nextTriggerTime.Format(utils.DATE_FORMAT)),
				Time:            nextTriggerTime.Format(utils.TIME_ONLY_FORMAT),
				ReminderText:    strings.TrimPrefix(reminderText, utils.REMINDER_PREFIX),
				InConstruction:  false,
				NextTriggerTime: nextTriggerTime.In(time.UTC).Format(utils.DIRECTUS_DATETIME_FORMAT),
			}
			err = reminder.Create()
			if err != nil {
				log.Fatal(err)
			}
			editedMessage := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				fmt.Sprintf("%v\n\nI will remind you again on %v", reminderText, nextTriggerTime.Format(utils.DATE_AND_TIME_FORMAT)),
			)
			if _, err := bot.Request(editedMessage); err != nil {
				log.Fatal(err)
			}
		case utils.RENEW_REMINDER_1D:
			nextTriggerTime := time.Now().In(tz).Add(24 * time.Hour)
			reminder := schemas.Reminder{
				Id:              uuid.New().String(),
				ChatId:          update.CallbackQuery.Message.Chat.ID,
				FromUserId:      update.CallbackQuery.From.ID,
				FileId:          "",
				Timezone:        chatSettings.Timezone,
				Frequency:       fmt.Sprintf("%v-%v", utils.REMINDER_ONCE, nextTriggerTime.Format(utils.DATE_FORMAT)),
				Time:            nextTriggerTime.Format(utils.TIME_ONLY_FORMAT),
				ReminderText:    strings.TrimPrefix(reminderText, utils.REMINDER_PREFIX),
				InConstruction:  false,
				NextTriggerTime: nextTriggerTime.In(time.UTC).Format(utils.DIRECTUS_DATETIME_FORMAT),
			}
			err = reminder.Create()
			if err != nil {
				log.Fatal(err)
			}
			editedMessage := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				fmt.Sprintf("%v\n\nI will remind you again on %v", reminderText, nextTriggerTime.Format(utils.DATE_AND_TIME_FORMAT)),
			)
			if _, err := bot.Request(editedMessage); err != nil {
				log.Fatal(err)
			}
		case utils.RENEW_REMINDER_CUSTOM:
			reminder := schemas.Reminder{
				Id:              uuid.New().String(),
				ChatId:          update.CallbackQuery.Message.Chat.ID,
				FromUserId:      update.CallbackQuery.From.ID,
				FileId:          "",
				Timezone:        chatSettings.Timezone,
				Frequency:       "",
				Time:            "",
				ReminderText:    strings.TrimPrefix(reminderText, utils.REMINDER_PREFIX),
				InConstruction:  true,
				NextTriggerTime: "",
			}
			err = reminder.Create()
			if err != nil {
				log.Fatal(err)
			}
			editedMessage := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				reminderText,
			)
			if _, err := bot.Request(editedMessage); err != nil {
				log.Fatal(err)
			}
			newMsg := tgbotapi.NewMessage(
				update.CallbackQuery.Message.Chat.ID,
				fmt.Sprintf("@%v enter reminder time in <HH>:<MM> format.", update.CallbackQuery.From.UserName),
			)
			newMsg.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: true,
				Selective:  true,
			}
			if _, err := bot.Request(newMsg); err != nil {
				log.Fatal(err)
			}
		case utils.RENEW_REMINDER_CANCEL:
			editedMessage := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				reminderText,
			)
			if _, err := bot.Request(editedMessage); err != nil {
				log.Fatal(err)
			}
		default:
			return
		}
	}
}
