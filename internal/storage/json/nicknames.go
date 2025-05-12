package json

import (
	"UserInteractiveBot/internal/models"
	"encoding/json"
	"os"
)

type JSONNicknameStorage struct {
	filePath string
}

func NewJSONNicknameStorage(filePath string) *JSONNicknameStorage {
	return &JSONNicknameStorage{filePath: filePath}
}

func (s *JSONNicknameStorage) Load() ([]models.UserNickname, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		emptyList := []models.UserNickname{}
		data, err := json.Marshal(emptyList)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(s.filePath, data, 0644); err != nil {
			return nil, err
		}
		return emptyList, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var list []models.UserNickname
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *JSONNicknameStorage) Save(list []models.UserNickname) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}
