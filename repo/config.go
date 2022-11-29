package repo

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

const DEFAULT_LANG string = "en"

type Config struct {
	Username string `json:"username"`
	Lang     string `json:"lang"`
	ClientID string `json:"clientID"`
}

type repoConfigFile struct {
	filename string
	fileMode fs.FileMode
	config   *Config
}

func NewRepoConfigFile(filename string) (*repoConfigFile, error) {
	config, err := getConfig(filename)

	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		config = &Config{}
	}

	return &repoConfigFile{
		filename: filename,
		fileMode: os.FileMode(0666),
		config:   config,
	}, nil
}

func (r *repoConfigFile) GetAnonymousUsername() string {
	return r.config.Username
}

func (r *repoConfigFile) GetClientID() (string, error) {
	if r.config.ClientID == "" {
		return "", errors.New("Missing ClientID in the config.json file")
	}
	return r.config.ClientID, nil
}

func (r *repoConfigFile) GetLang() string {
	if r.config.Lang == "" {
		return DEFAULT_LANG
	}
	return r.config.Lang
}

func getConfig(filename string) (*Config, error) {
	buff, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = json.Unmarshal(buff, c)

	return c, err
}

func (r *repoConfigFile) save() error {
	buff, err := json.Marshal(r.config)

	if err != nil {
		return err
	}

	return os.WriteFile(r.filename, buff, r.fileMode)
}

func (r *repoConfigFile) SaveAnonymousUsername(username string) error {
	r.config.Username = username
	return r.save()
}

func (r *repoConfigFile) SaveLang(lang string) error {
	r.config.Lang = lang
	return r.save()
}
