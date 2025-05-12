package json

import (
	"UserInteractiveBot/internal/storage"
)

type JSONStorage struct {
	users     *JSONUserStorage
	admins    *JSONUserStorage
	nicknames *JSONNicknameStorage
	posts     *JSONPostStorage
}

func NewJSONStorage(usersPath, adminsPath, nicknamesPath, postsPath string) *JSONStorage {
	return &JSONStorage{
		users:     NewJSONUserStorage(usersPath),
		admins:    NewJSONUserStorage(adminsPath),
		nicknames: NewJSONNicknameStorage(nicknamesPath),
		posts:     NewJSONPostStorage(postsPath),
	}
}

func (s *JSONStorage) Users() storage.UserStorage {
	return s.users
}

func (s *JSONStorage) Admins() storage.UserStorage {
	return s.admins
}

func (s *JSONStorage) Nicknames() storage.NicknameStorage {
	return s.nicknames
}

func (s *JSONStorage) Posts() storage.PostStorage {
	return s.posts
}
