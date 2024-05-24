package utils

var (
	LogLevel      = "info"
	DirectusHost  = "http://localhost:8055"
	DirectusToken = "directus-access-token"
	BotToken      = "my-bot-token"
)

const HELP_MESSAGE string = `This bot lets you set reminders! The following commands are available:
/remind sets a reminder.
/list displays all the reminders in the current chat.
/settings to set timezone.


Note that all reminders set on this bot can be accessed by the user hosting this bot. Do not set any reminders that contain any sort of private information.`
const SUPPORT_MESSAGE string = `My source code is hosted on https://github.com/Jason-CKY/telegram-reminderbot. Consider 
Post any issues with this bot on the github link, and feel free to contribute to the source code with a pull request.`
const DEFAULT_SETTINGS_MESSAGE = "The default timezone for the bot is Asia/Singapore (GMT +8). Type /settings for more information on how to change the timezone. "
const REMINDER_BUILDER_MESSAGE string = `Please enter reminder text. This bot allows for image reminders as well. Just attach an image and put your reminder text as the caption.`
const CANCEL_MESSAGE string = `🚫 Cancel`
const CANCEL_OPERATION_MESSAGE string = `Operation cancelled.`
const DEFAULT_TIMEZONE = "Asia/Singapore"

const REMINDER_ONCE = "Once"
const REMINDER_DAILY = "Daily"
const REMINDER_WEEKLY = "Weekly"
const REMINDER_MONTHLY = "Monthly"
const REMINDER_YEARLY = "Yearly"

var DAY_OF_WEEK = map[string]int{
	"Sunday":    0,
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
}

const CALLBACK_NO_ACTION = "n"
const CALLBACK_GOTO = "g"
const CALLBACK_SELECT = "s"
const CALLBACK_DELETE = "d"
const CALLBACK_SHOW_IMAGE = "p"

const CALLBACK_CALENDAR_STEP_YEAR = "y"
const CALLBACK_CALENDAR_STEP_MONTH = "m"
const CALLBACK_CALENDAR_STEP_DAY = "d"

const CALLBACK_CALENDAR_SELECT_YEAR = "Select year"
const CALLBACK_CALENDAR_SELECT_MONTH = "Select month"
const CALLBACK_CALENDAR_SELECT_DAY = "Select day"

const CALENDAR_YEAR_NUM_ROWS = 2
const CALENDAR_YEAR_NUM_COLS = 2
const CALENDAR_MONTH_NUM_ROWS = 4
const CALENDAR_MONTH_NUM_COLS = 3

// Day calendar doesn't need a hard-coded number of rows as they need to be under the correct day columns so rows is not fixed
const CALENDAR_DAY_NUM_COLS = 7

// format dates in YYYY/MM/DD
const DATE_FORMAT = "2006/01/02"
const PRETTY_DATE_FORMAT = "Mon, 02 Jan 2006"
const PRETTY_DATE_FORMAT_WITHOUT_YEAR = "02 Jan"
const TIME_ONLY_FORMAT = "15:04"
const DATE_AND_TIME_FORMAT = "Mon, 02 Jan 2006 15:04:05"
const DATE_AND_TIME_FORMAT_WITHOUT_YEAR = "02 Jan 15:04:05"
const DIRECTUS_DATETIME_FORMAT = "2006-01-02T15:04:05"

const REMINDER_PREFIX = "🗓"
const REMINDER_PHOTO_PREFIX = "🖼"
const RENEW_REMINDER_15M = "renew_15m"
const RENEW_REMINDER_30M = "renew_30m"
const RENEW_REMINDER_1H = "renew_1h"
const RENEW_REMINDER_3H = "renew_3h"
const RENEW_REMINDER_1D = "renew_1d"
const RENEW_REMINDER_CUSTOM = "renew_time"
const RENEW_REMINDER_CANCEL = "renew_cancel"
const RENEW_REMINDER_TEXT = "\n\nRemind me again in:"

const SETTINGS_CHANGE_TIMEZONE = "🕐 Change time zone"
const CHANGE_TIMEZONE_MESSAGE = "Please type the timezone that you want to change to. For a list of all supported timezones, please click click <a href=\"https://timeapi.io/documentation/iana-timezones\">here</a>"
const INVALID_TIMEZONE_MESSAGE = "Invalid timezone.\n\nFor a list of all supported timezones, please click <a href=\"https://gist.github.com/heyalexej/8bf688fd67d7199be4a1682b3eec7568\">here</a>"

// list reminder settings
const MAX_REMINDERS_PER_PAGE = 5
const NO_REMINDERS_MESSAGE = "There are no reminders on current page, try to open another page or request list again."
