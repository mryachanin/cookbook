package config

import (
  "encoding/json"
  "log"
  "os"
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

func LoadConfiguration(file string) *Config {
  var config *Config
  configFile, err := os.Open(file)
  defer configFile.Close()
  if err != nil {
    log.Fatal(err.Error())
  }
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&config)
  return config
}
