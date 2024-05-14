package core

import (
	"time"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
	log "github.com/sirupsen/logrus"
)

func ScheduledReminderTrigger() {
	for {
		log.Info("Checking db for reminder triggers..")
		dueReminders, err := schemas.GetDueReminders()
		if err != nil {
			panic(err)
		}
		log.Info(len(dueReminders))
		time.Sleep(2 * time.Second)
	}
}
