package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ConfigDir  string
	ConfigFile string
	Settings   map[string]interface{}
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
		Settings:   make(map[string]interface{}),
	}

	// Set default values
	defaults := map[string]interface{}{
		"db_path":      "go-time.db", // The relative path to the DB
		"command_mode": "cli",        // Default command mode
	}

	if _, err := os.Stat(configFile); err == nil {
		c.Load()
	} else {
		c.Settings = defaults
		c.Save()
	}

	return c
}

func (c *Config) Load() error {
	data, err := os.ReadFile(c.ConfigFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c.Settings)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c.Settings, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(c.ConfigFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Set(key string, value interface{}) {
	c.Settings[key] = value
	c.Save()
}

func (c *Config) Get(key string, defaultValue interface{}) interface{} {
	if value, ok := c.Settings[key]; ok {
		return value
	}
	return defaultValue
}

func (c *Config) Remove(key string) {
	delete(c.Settings, key)
	c.Save()
}

// Example usage:
func main() {
	config := New("go-time", "config.json")
	config.Set("db_path", "go-time.db")
	fmt.Println(config.Get("db_path", "default.db"))
	config.Remove("db_path")
	fmt.Println(config.Get("db_path", "default.db"))
}
