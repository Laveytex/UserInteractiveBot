package models

import "time"

type Post struct {
	ID          string    `json:"id"`
	AdminChatID int64     `json:"admin_chat_id"`
	Text        string    `json:"text"`
	PhotoIDs    []string  `json:"photo_ids"`
	PublishTime time.Time `json:"publish_time"`
	Published   bool      `json:"published"`
}

type CreatingPost struct {
	AdminChatID int64
	Text        string
	PhotoIDs    []string
	State       AdminState
}

type AdminState string

const (
	StateIdle        AdminState = "idle"
	StateAddingText  AdminState = "adding_text"
	StateAddingPhoto AdminState = "adding_photo"
	StateSettingTime AdminState = "setting_time"
)
