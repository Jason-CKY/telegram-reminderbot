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
