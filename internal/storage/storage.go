package storage

import (
	"UserInteractiveBot/internal/models"
)

type Storage interface {
	Users() UserStorage
	Admins() UserStorage
	Nicknames() NicknameStorage
	Posts() PostStorage
}

type UserStorage interface {
	Load() ([]int64, error)
	Save([]int64) error
}

type NicknameStorage interface {
	Load() ([]models.UserNickname, error)
	Save([]models.UserNickname) error
}

type PostStorage interface {
	Load() ([]models.Post, error)
	Save([]models.Post) error
}
