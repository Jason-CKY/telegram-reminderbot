package core

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
