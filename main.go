package main

import (
	"flag"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/core"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/handler"
	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()

	if err != nil {
		log.Infof("Error loading .env file: %v\nUsing environment variables instead...", err)
	}

	flag.StringVar(&utils.LogLevel, "log-level", utils.LookupEnvOrString("LOG_LEVEL", utils.LogLevel), "Logging level for the server")
	flag.StringVar(&utils.DirectusHost, "directus-host", utils.LookupEnvOrString("DIRECTUS_HOST", utils.DirectusHost), "Hostname for directus server")
	flag.StringVar(&utils.DirectusToken, "directus-token", utils.LookupEnvOrString("DIRECTUS_TOKEN", utils.DirectusToken), "Access token for directus")
	flag.StringVar(&utils.BotToken, "bot-token", utils.LookupEnvOrString("TELEGRAM_BOT_TOKEN", utils.BotToken), "Bot token for telegram bot")

	flag.Parse()

	// setup logrus
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	logLevel, _ := log.ParseLevel(utils.LogLevel)
	log.SetLevel(logLevel)

	log.Infof("connecting to directus at: %v", utils.DirectusHost)

	bot, err := tgbotapi.NewBotAPI(utils.BotToken)
	if err != nil {
		panic(fmt.Errorf("telegram connection error: %v", err))
	}
	bot.Debug = utils.LogLevel == "debug"
	log.Infof("Authorized on account %s", bot.Self.UserName)

	if err != nil {
		panic(err)
	}

	go core.ScheduledReminderTrigger(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		handler.HandleUpdate(&update, bot)
	}

}
