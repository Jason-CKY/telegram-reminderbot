package schemas

type Reminder struct {
	Id         string `json:"id"`
	ChatId     int64  `json:"chat_id"`
	FromUserId int64  `json:"from_user_id"`
	FileId     string `json:"file_id"`
	Timezone   string `json:"timezone"`
	Frequency  string `json:"frequency"`
	Time       string `json:"time"`
}
