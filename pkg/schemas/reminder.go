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
	log "github.com/sirupsen/logrus"
)

type Reminder struct {
	Id              string `json:"id"`
	ChatId          int64  `json:"chat_id"`
	FromUserId      int64  `json:"from_user_id"`
	FileId          string `json:"file_id"`
	Frequency       string `json:"frequency"`
	Time            string `json:"time"`
	ReminderText    string `json:"reminder_text"`
	InConstruction  bool   `json:"in_construction"`
	NextTriggerTime string `json:"next_trigger_time,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
func (r Reminder) MarshalJSON() ([]byte, error) {
	type Alias Reminder // Create an alias to avoid recursion

	aux := &struct {
		ChatId     string `json:"chat_id"`
		FromUserId string `json:"from_user_id"`
		*Alias
	}{
		ChatId:     strconv.FormatInt(r.ChatId, 10),
		FromUserId: strconv.FormatInt(r.FromUserId, 10),
		Alias:      (*Alias)(&r),
	}
	return json.Marshal(aux)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Reminder) UnmarshalJSON(data []byte) error {
	type Alias Reminder // Create an alias to avoid recursion

	aux := &struct {
		ChatId     interface{} `json:"chat_id"`
		FromUserId interface{} `json:"from_user_id"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		log.Error("Error unmarshaling reminder: ", err)
		return err
	}

	// Handle chat_id as string or number
	switch v := aux.ChatId.(type) {
	case string:
		chatId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Error("Error parsing chat_id: ", err)
			return err
		}
		r.ChatId = chatId
	case float64:
		r.ChatId = int64(v)
	default:
		return fmt.Errorf("unexpected type for chat_id: %T", v)
	}

	// Handle from_user_id as string or number
	switch v := aux.FromUserId.(type) {
	case string:
		fromUserId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Error("Error parsing from_user_id: ", err)
			return err
		}
		r.FromUserId = fromUserId
	case float64:
		r.FromUserId = int64(v)
	default:
		return fmt.Errorf("unexpected type for from_user_id: %T", v)
	}

	return nil
}

func (reminder Reminder) Create() error {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder", utils.DirectusHost)
	reqBody, _ := json.Marshal(reminder)
	req, httpErr := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error inserting reminder to directus: %v", string(body))
	}

	return nil
}

func (reminder Reminder) Update() error {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder/%v", utils.DirectusHost, reminder.Id)
	reqBody, _ := json.Marshal(reminder)
	req, httpErr := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error updating reminder to directus: %v", string(body))
	}

	return nil
}

func (reminder Reminder) Delete() error {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder/%v", utils.DirectusHost, reminder.Id)
	req, httpErr := http.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 204 {
		return fmt.Errorf("error deleting reminder in directus: %v", string(body))
	}
	return nil
}

func (reminder Reminder) DeleteReminderInConstruction() error {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder", utils.DirectusHost)
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 && res.StatusCode != 204 {
		return fmt.Errorf("error deleting reminders in construction: %v", string(body))
	}

	return nil
}

func (reminder Reminder) CalculateNextTriggerTime(chatSettings *ChatSettings) (time.Time, error) {
	// calculate the next trigger time, in the user's timezone
	tz, _ := time.LoadLocation(chatSettings.Timezone)
	frequencyText := strings.Split(reminder.Frequency, "-")
	frequency := frequencyText[0]
	if reminder.Frequency == utils.REMINDER_DAILY {
		frequency = utils.REMINDER_DAILY
	}
	switch {
	case frequency == utils.REMINDER_ONCE:
		triggerTime, err := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%v %v", frequencyText[1], reminder.Time), tz)
		if err != nil {
			return time.Now(), err
		}
		return triggerTime.In(time.UTC), nil
	case frequency == utils.REMINDER_DAILY:
		currentTime := time.Now().UTC()
		reminderHour, reminderMinute := utils.ParseReminderTime(reminder.Time)
		triggerTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), reminderHour, reminderMinute, 0, 0, tz).In(time.UTC)
		if currentTime.After(triggerTime) {
			return triggerTime.Add(24 * time.Hour), nil
		}
		return triggerTime, nil
	case frequency == utils.REMINDER_WEEKLY:
		reminderWeekday, _ := strconv.Atoi(frequencyText[1])
		currentTime := time.Now().In(tz)
		reminderHour, reminderMinute := utils.ParseReminderTime(reminder.Time)
		triggerTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), reminderHour, reminderMinute, 0, 0, tz)
		for reminderWeekday != int(triggerTime.In(tz).Weekday()) {
			triggerTime = triggerTime.Add(24 * time.Hour)
		}
		if currentTime.After(triggerTime) {
			triggerTime = triggerTime.Add(7 * 24 * time.Hour)
		}

		return triggerTime.In(time.UTC), nil
	case frequency == utils.REMINDER_MONTHLY:
		reminderDay, _ := strconv.Atoi(frequencyText[1])
		currentTime := time.Now().In(tz)
		reminderHour, reminderMinute := utils.ParseReminderTime(reminder.Time)
		triggerTime := time.Date(currentTime.Year(), currentTime.Month(), reminderDay, reminderHour, reminderMinute, 0, 0, tz)

		if currentTime.After(triggerTime) {
			return time.Date(currentTime.Year(), currentTime.Month()+1, reminderDay, reminderHour, reminderMinute, 0, 0, tz).In(time.UTC), nil
		}
		return triggerTime.In(time.UTC), nil
	case frequency == utils.REMINDER_YEARLY:
		t, err := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%v %v", frequencyText[1], reminder.Time), tz)
		if err != nil {
			return time.Now(), err
		}
		currentTime := time.Now().In(tz)
		triggerTime := time.Date(currentTime.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, tz)
		if currentTime.After(triggerTime) {
			return time.Date(currentTime.Year()+1, t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, tz).In(time.UTC), nil
		}
		return triggerTime.In(time.UTC), nil
	default:
		return time.Now(), errors.New("invalid frequency")
	}
}

func GetReminderInConstruction(chatId int64, fromUserId int64) (*Reminder, error) {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder", utils.DirectusHost)
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return nil, httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return nil, httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error searching for reminder in directus: %v", string(body))
	}
	log.Debug(string(body))
	var reminderResponse map[string][]Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		log.Error("Error unmarshaling reminder response: ", jsonErr)
		return nil, jsonErr
	}
	log.Debug(reminderResponse)

	if len(reminderResponse["data"]) == 0 {
		return nil, nil
	}

	return &reminderResponse["data"][0], nil
}

func GetReminderById(Id string) (*Reminder, error) {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder", utils.DirectusHost)
	reqBody := []byte(fmt.Sprintf(`{
		"query": {
			"filter": {
				"id": {
					"_eq": "%v"
				}
			}
		}
	}`, Id))
	req, httpErr := http.NewRequest("SEARCH", endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return nil, httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return nil, httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error searching for reminder in directus: %v", string(body))
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

func GetRemindersByChatId(chatId int64) ([]Reminder, error) {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder", utils.DirectusHost)
	reqBody := []byte(fmt.Sprintf(`{
		"query": {
			"filter": {
				"chat_id": {
					"_eq": "%v"
				},
				"in_construction": {
					"_eq": false
				}
			}
		}
	}`, chatId))
	req, httpErr := http.NewRequest("SEARCH", endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return nil, httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return nil, httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error searching for reminder in directus: %v", string(body))
	}
	var reminderResponse map[string][]Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return nil, jsonErr
	}

	return reminderResponse["data"], nil
}

func GetDueReminders() ([]Reminder, error) {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder", utils.DirectusHost)
	reqBody := []byte(fmt.Sprintf(`{
		"query": {
			"filter": {
				"_and": [
					{
						"in_construction": {
							"_eq": false
						}
					},
					{
						"next_trigger_time": {
							"_lt": "%v"
						}
					}
				]
			}
		}
	}`, time.Now().UTC().Format(utils.DIRECTUS_DATETIME_FORMAT)))
	req, httpErr := http.NewRequest("SEARCH", endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
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
		return nil, fmt.Errorf("error searching for reminder in directus: %v", string(body))
	}
	log.Debug(string(body))
	var reminderResponse map[string][]Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return nil, jsonErr
	}

	return reminderResponse["data"], nil
}

func ListChatReminders(chatId int64) ([]Reminder, error) {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder", utils.DirectusHost)
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
						"in_construction": {
							"_eq": false
						}
					}
				]
			},
			"sort": "date_created"
		}
	}`, chatId))
	req, httpErr := http.NewRequest("SEARCH", endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return nil, httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return nil, httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error searching for reminders in directus: %v", string(body))
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
	return reminderResponse["data"], nil
}

func MigrateReminderChatId(fromChatId int64, toChatId int64) error {
	endpoint := fmt.Sprintf("%v/items/reminderbot_reminder", utils.DirectusHost)
	reqBody := []byte(fmt.Sprintf(`{
		"query": {
			"filter": {
				"chat_id": {
					"_eq": "%v"
				}
			}
		},
		"data": {
			"chat_id": "%v"
		}
	}`, fromChatId, toChatId))
	req, httpErr := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", utils.DirectusToken))
	if httpErr != nil {
		return httpErr
	}
	client := &http.Client{}
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error updating reminders in directus: %v", string(body))
	}
	var reminderResponse map[string][]Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return jsonErr
	}
	return nil
}
