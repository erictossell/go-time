package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ConfigDir  string
	ConfigFile string
	Settings   AppConfig
}

type AppConfig struct {
	DBPath      string `toml:"db_path"`
	CommandMode string `toml:"command_mode"`
}

func New(appname string, filename string) *Config {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", appname)
	configFile := filepath.Join(configDir, filename)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, os.ModePerm)
	}

	c := &Config{
		ConfigDir:  configDir,
		ConfigFile: configFile,
		Settings: AppConfig{
			DBPath:      "go-time.db",
			CommandMode: "cli",
		},
	}

	if _, err := os.Stat(configFile); err == nil {
		c.Load()
	} else {
		c.SaveWithComments()
	}

	return c
}

func (c *Config) Load() error {
	_, err := toml.DecodeFile(c.ConfigFile, &c.Settings)
	return err
}

func (c *Config) Save() error {
	f, err := os.Create(c.ConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(c.Settings)
}

func (c *Config) SaveWithComments() error {
	configWithComments := `# Path to the SQLite database file
db_path = "` + c.Settings.DBPath + `"

# Mode in which the application runs (cli, tui, help)
command_mode = "` + c.Settings.CommandMode + `"
`

	err := os.WriteFile(c.ConfigFile, []byte(configWithComments), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Set(key string, value interface{}) {
	switch key {
	case "db_path":
		c.Settings.DBPath = value.(string)
	case "command_mode":
		c.Settings.CommandMode = value.(string)
	}
	c.Save()
}

func (c *Config) Get(key string, defaultValue interface{}) interface{} {
	switch key {
	case "db_path":
		if c.Settings.DBPath != "" {
			return c.Settings.DBPath
		}
	case "command_mode":
		if c.Settings.CommandMode != "" {
			return c.Settings.CommandMode
		}
	}
	return defaultValue
}

func (c *Config) Remove(key string) {
	switch key {
	case "db_path":
		c.Settings.DBPath = ""
	case "command_mode":
		c.Settings.CommandMode = ""
	}
	c.Save()
}
