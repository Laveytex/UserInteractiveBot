package json

import (
	"UserInteractiveBot/internal/models"
	"encoding/json"
	"os"
)

type JSONPostStorage struct {
	filePath string
}

func NewJSONPostStorage(filePath string) *JSONPostStorage {
	return &JSONPostStorage{filePath: filePath}
}

func (s *JSONPostStorage) Load() ([]models.Post, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		emptyList := []models.Post{}
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

	var list []models.Post
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *JSONPostStorage) Save(list []models.Post) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}
