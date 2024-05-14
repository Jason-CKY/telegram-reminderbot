package utils

var (
	LogLevel     = "info"
	DirectusHost = "http://localhost:8055"
	WebPort      = 8080
	BotToken     = "my-bot-token"
)

const HELP_MESSAGE string = `This bot lets you set reminders! The following commands are available:
/remind sets a reminder.
/list displays all the reminders in the current chat.
/settings to set timezone.


Note that all reminders set on this bot can be accessed by the user hosting this bot. Do not set any reminders that contain any sort of private information.`

const SUPPORT_MESSAGE string = `My source code is hosted on https://github.com/Jason-CKY/telegram-reminderbot. Consider 
Post any issues with this bot on the github link, and feel free to contribute to the source code with a pull request.`

const REMINDER_BUILDER_MESSAGE string = `Please enter reminder text. This bot allows for image reminders as well. Just attach an image and put your reminder text as the caption.`

const CANCEL_MESSAGE string = `🚫 Cancel`

const CANCEL_OPERATION_MESSAGE string = `Operation cancelled.`

const REMINDER_ONCE = "Once"
const REMINDER_DAILY = "Daily"
const REMINDER_WEEKLY = "Weekly"
const REMINDER_MONTHLY = "Monthly"
const REMINDER_YEARLY = "Yearly"

const MONDAY = "Monday"
const TUESDAY = "Tuesday"
const WEDNESDAY = "Wednesday"
const THURSDAY = "Thursday"
const FRIDAY = "Friday"
const SATURDAY = "Saturday"
const SUNDAY = "Sunday"

var DAY_OF_WEEK = map[string]int{
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
	"Sunday":    7,
}

const CALLBACK_NO_ACTION = "n"
const CALLBACK_GOTO = "g"
const CALLBACK_SELECT = "s"

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
