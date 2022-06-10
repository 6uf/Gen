package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

func (s *Config) ToJson() (Data []byte, err error) {
	return json.MarshalIndent(s, "", "  ")
}

func (config *Config) SaveConfig() {
	if Json, err := config.ToJson(); err == nil {
		WriteFile("config.json", string(Json))
	}
}

func (s *Config) LoadState() {
	data, err := ReadFile("config.json")
	if err != nil {
		s.Key = uuid.NewString()[0:32]
		s.LoadFromFile()
		s.SaveConfig()
		return
	}

	json.Unmarshal([]byte(data), s)
	s.LoadFromFile()
}

func (c *Config) LoadFromFile() {
	// Load a config file

	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		jsonFile, _ = os.Create("config.json")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func WriteFile(path string, content string) {
	ioutil.WriteFile(path, []byte(content), 0644)
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
