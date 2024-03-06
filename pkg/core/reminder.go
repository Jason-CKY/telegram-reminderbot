package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/schemas"
)

func CreateReminder(reminder *schemas.Reminder) error {
	endpoint := fmt.Sprintf("%v/items/reminder", DirectusHost)
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
	var reminderResponse map[string]schemas.Reminder
	jsonErr := json.Unmarshal(body, &reminderResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}
