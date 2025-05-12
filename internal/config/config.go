package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram struct {
		TokenEnv string `yaml:"token_env"`
		Debug    bool   `yaml:"debug"`
	} `yaml:"telegram"`
	Auth struct {
		SecretCode string `yaml:"secret_code"`
		AdminCode  string `yaml:"admin_code"`
	} `yaml:"auth"`
	Storage struct {
		UsersFile     string `yaml:"users_file"`
		AdminsFile    string `yaml:"admins_file"`
		NicknamesFile string `yaml:"nicknames_file"`
		PostsFile     string `yaml:"posts_file"`
	} `yaml:"storage"`
	Limits struct {
		MaxCaptionLength int `yaml:"max_caption_length"`
		MaxMediaCount    int `yaml:"max_media_count"`
	} `yaml:"limits"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
