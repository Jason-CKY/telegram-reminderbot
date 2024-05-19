package schemas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jason-CKY/telegram-reminderbot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChatSettings struct {
	ChatId   int64  `json:"chat_id"`
	Timezone string `json:"timezone"`
	Updating bool   `json:"updating"`
}

func (chatSettings ChatSettings) Create() error {
	endpoint := fmt.Sprintf("%v/items/chat_settings", utils.DirectusHost)
	reqBody, _ := json.Marshal(chatSettings)
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
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error inserting chat settings to directus: %v", string(body))
	}

	return nil
}

func (chatSettings ChatSettings) Update() error {
	endpoint := fmt.Sprintf("%v/items/chat_settings/%v", utils.DirectusHost, chatSettings.ChatId)
	reqBody, _ := json.Marshal(chatSettings)
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
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error updating chat settings to directus: %v", string(body))
	}

	return nil
}

func GetChatSettings(chatId int64) (*ChatSettings, error) {
	endpoint := fmt.Sprintf("%v/items/chat_settings", utils.DirectusHost)
	reqBody := []byte(fmt.Sprintf(`{
		"query": {
			"filter": {
				"chat_id": {
					"_eq": "%v"
				}
			}
		}
	}`, chatId))
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
		return nil, fmt.Errorf("error getting chat settings in directus: %v", string(body))
	}
	var chatSettingsResponse map[string][]ChatSettings
	jsonErr := json.Unmarshal(body, &chatSettingsResponse)
	// error handling for json unmarshaling
	if jsonErr != nil {
		return nil, jsonErr
	}

	if len(chatSettingsResponse["data"]) == 0 {
		return nil, nil
	}

	return &chatSettingsResponse["data"][0], nil
}

func InsertChatSettingsIfNotPresent(chatId int64, bot *tgbotapi.BotAPI) (*ChatSettings, error) {
	chatSettings, err := GetChatSettings(chatId)
	if err != nil {
		return nil, err
	}
	if chatSettings == nil {
		chatSettings = &ChatSettings{
			ChatId:   chatId,
			Timezone: utils.DEFAULT_TIMEZONE,
			Updating: false,
		}
		err := chatSettings.Create()
		if err != nil {
			return nil, err
		}
		msg := tgbotapi.NewMessage(
			chatId,
			utils.DEFAULT_SETTINGS_MESSAGE,
		)
		if _, err := bot.Request(msg); err != nil {
			return nil, err
		}
	}
	return chatSettings, nil
}
