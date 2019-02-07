package config

import (
  "encoding/json"
  "log"
  "os"
)

const (
  configPath = "./config.json"
)

type Config struct {
  Host string       `json:"host"`
  Port int          `json:"port"`
  Database struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"username"`
    Password string `json:"password"`
  }                 `json:"database"`
}

func LoadConfiguration() *Config {
  var config *Config
  configFile, err := os.Open(configPath)
  defer configFile.Close()
  if err != nil {
    log.Fatal(err.Error())
  }
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&config)
  return config
}
