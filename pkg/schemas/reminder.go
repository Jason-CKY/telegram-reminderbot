package schemas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/core"
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

// func getReminderBy

func (reminder Reminder) Create() error {
	endpoint := fmt.Sprintf("%v/items/reminder", core.DirectusHost)
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
