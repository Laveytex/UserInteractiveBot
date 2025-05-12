package json

import (
	"encoding/json"
	"os"
)

type JSONUserStorage struct {
	filePath string
}

func NewJSONUserStorage(filePath string) *JSONUserStorage {
	return &JSONUserStorage{filePath: filePath}
}

func (s *JSONUserStorage) Load() ([]int64, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		emptyList := []int64{}
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

	var list []int64
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *JSONUserStorage) Save(list []int64) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}
