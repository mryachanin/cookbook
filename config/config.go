package config

import (
  "encoding/json"
  "log"
  "os"
)

type Config struct {
  Host string `json:"host"`
  Port int `json:"port"`
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
