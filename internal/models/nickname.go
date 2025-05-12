package models

type UserNickname struct {
	ChatID   int64  `json:"chat_id"`
	Nickname string `json:"nickname"`
}
