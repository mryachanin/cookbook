// Defines the config path for db operations
package ops

import "os"

var configPath = getConfigPath()

func getConfigPath() string {
  if path := os.Getenv("COOKBOOK_CONFIG"); path != "" {
    return path
  }
  if _, err := os.Stat("config_local.json"); err == nil {
    return "config_local.json"
  }
  if _, err := os.Stat("../../config_json"); err == nil {
    return "../../config_json"
  }
  return "config_json"
}