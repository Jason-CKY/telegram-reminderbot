package schemas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
)

type Reminder struct {
	Id              string `json:"id"`
	ChatId          int64  `json:"chat_id"`
	FromUserId      int64  `json:"from_user_id"`
	FileId          string `json:"file_id"`
	Timezone        string `json:"timezone"`
	Frequency       string `json:"frequency"`
	Time            string `json:"time"`
	ReminderText    string `json:"reminder_text"`
	InConstruction  bool   `json:"in_construction"`
	NextTriggerTime string `json:"next_trigger_time"`
}

func (reminder Reminder) Create() error {
	endpoint := fmt.Sprintf("%v/items/reminder", utils.DirectusHost)
	reqBody, _ := json.Marshal(reminder)
	req, httpErr := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return fmt.Errorf("error inserting reminder to directus: %v", res.Status)
	}
	var reminderResponse map[string]Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

func (reminder Reminder) Update() error {
	endpoint := fmt.Sprintf("%v/items/reminder/%v", utils.DirectusHost, reminder.Id)
	reqBody, _ := json.Marshal(reminder)
	req, httpErr := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return fmt.Errorf("error updating reminder to directus: %v", res.Status)
	}
	var reminderResponse map[string]Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

func (reminder Reminder) DeleteById() error {
	endpoint := fmt.Sprintf("%v/items/reminder/%v", utils.DirectusHost, reminder.Id)
	req, httpErr := http.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return fmt.Errorf("error deleting reminder in directus: %v", res.Status)
	}
	var reminderResponse map[string]Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

func (reminder Reminder) DeleteReminderInConstruction() error {
	// delete all reminders that has from_user_id == reminder.FromUserId and ChatId == reminder.ChatId
	endpoint := fmt.Sprintf("%v/items/reminder", utils.DirectusHost)
	reqBody := []byte(fmt.Sprintf(`{
		"query": {
			"filter": {
				"_and": [
					{
						"chat_id": {
							"_eq": "%v"
						}
					},
					{
						"from_user_id": {
							"_eq": "%v"
						}
					},
					{
						"in_construction": {
							"_eq": true
						}
					}
				]
			}
		}
	}`, reminder.ChatId, reminder.FromUserId))

	req, httpErr := http.NewRequest(http.MethodDelete, endpoint, bytes.NewBuffer(reqBody))

	req.Header.Set("Content-Type", "application/json")
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	defer res.Body.Close()
	if res.StatusCode != 200 && res.StatusCode != 204 {
		return fmt.Errorf("error deleting reminders in construction: %v", res.Status)
	}

	return nil
}

func (reminder Reminder) CalculateNextTriggerTime() (time.Time, error) {
	// calculate the next trigger time, in the user's timezone
	frequencyText := strings.Split(reminder.Frequency, "-")
	frequency := frequencyText[0]
	switch {
	case frequency == utils.REMINDER_ONCE:
		tz, _ := time.LoadLocation(reminder.Timezone)
		// reminderTime stored in db is in UTC, while the date string is in user's timezone, so we need to correct that
		timezoneAdjustedTimeString := utils.ConvertReminderTimeFromUTCToUserTimezone(reminder.Time, tz)
		triggerTime, err := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%v %v", frequencyText[1], timezoneAdjustedTimeString), tz)
		if err != nil {
			return time.Now(), err
		}
		return triggerTime.In(time.UTC), nil
	case frequency == utils.REMINDER_DAILY:
		currentTime := time.Now().UTC()
		reminderHour, reminderMinute := utils.ParseReminderTime(reminder.Time)
		triggerTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), reminderHour, reminderMinute, 0, 0, time.UTC)
		if currentTime.After(triggerTime) {
			return triggerTime.Add(24 * time.Hour), nil
		}
		return triggerTime, nil
	case frequency == utils.REMINDER_WEEKLY:
		reminderWeekday, _ := strconv.Atoi(frequencyText[1])
		currentTime := time.Now().UTC()
		reminderHour, reminderMinute := utils.ParseReminderTime(reminder.Time)
		triggerTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), reminderHour, reminderMinute, 0, 0, time.UTC)

		reminderTz, _ := time.LoadLocation(reminder.Timezone)
		if reminderWeekday != int(triggerTime.In(reminderTz).Weekday()) || currentTime.After(triggerTime) {
			return triggerTime.Add(7 * 24 * time.Hour), nil
		}
		return triggerTime, nil
	case frequency == utils.REMINDER_MONTHLY:
		reminderDay, _ := strconv.Atoi(frequencyText[1])
		currentTime := time.Now().UTC()
		reminderHour, reminderMinute := utils.ParseReminderTime(reminder.Time)
		triggerTime := time.Date(currentTime.Year(), currentTime.Month(), reminderDay, reminderHour, reminderMinute, 0, 0, time.UTC)

		if currentTime.After(triggerTime) {
			return time.Date(currentTime.Year(), currentTime.Month()+1, reminderDay, reminderHour, reminderMinute, 0, 0, time.UTC), nil
		}
		return triggerTime, nil
	case frequency == utils.REMINDER_YEARLY:
		t, err := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%v %v", frequencyText[1], reminder.Time), time.UTC)
		if err != nil {
			return time.Now(), err
		}
		currentTime := time.Now().UTC()
		triggerTime := time.Date(currentTime.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
		if currentTime.After(triggerTime) {
			return time.Date(currentTime.Year()+1, t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC), nil
		}
		return triggerTime, nil
	default:
		return time.Now(), errors.New("invalid frequency")
	}
}

func GetReminderInConstruction(chatId int64, fromUserId int64) (*Reminder, error) {
	endpoint := fmt.Sprintf("%v/items/reminder", utils.DirectusHost)
	reqBody := []byte(fmt.Sprintf(`{
		"query": {
			"filter": {
				"_and": [
					{
						"chat_id": {
							"_eq": "%v"
						}
					},
					{
						"from_user_id": {
							"_eq": "%v"
						}
					},
					{
						"in_construction": {
							"_eq": true
						}
					}
				]
			}
		}
	}`, chatId, fromUserId))
	req, httpErr := http.NewRequest("SEARCH", endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if httpErr != nil {
		return nil, httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return nil, httpErr
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error searching for reminder in directus: %v", res.Status)
	}
	var reminderResponse map[string][]Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return nil, jsonErr
	}

	if len(reminderResponse["data"]) == 0 {
		return nil, nil
	}

	return &reminderResponse["data"][0], nil
}

func GetDueReminders() ([]Reminder, error) {
	dueReminders := []Reminder{}
	return dueReminders, nil
}
