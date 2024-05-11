package schemas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
)

type Reminder struct {
	Id             string `json:"id"`
	ChatId         int64  `json:"chat_id"`
	FromUserId     int64  `json:"from_user_id"`
	FileId         string `json:"file_id"`
	Timezone       string `json:"timezone"`
	Frequency      string `json:"frequency"`
	Time           string `json:"time"`
	ReminderText   string `json:"reminder_text"`
	InConstruction bool   `json:"in_construction"`
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
