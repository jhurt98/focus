package config
import (
    "fmt"
    "os"
    "encoding/json"
)

type Config struct {
    port int
    time int
}

var DefaultPort = 20002
var DefaultTime = -1

var configJsonFile = "./config.json"

func GetConfig() Config {
    // read config json
    // set port
    // set time
    // return Config struct
    config := &Config{}
    config.initConfig()
    fmt.Println("config", config)
    return *config
}

func (config *Config) initConfig() {
    readConfigJson(config)
}

func readConfigJson(config *Config) {
    jsonBytes, err := os.ReadFile(configJsonFile)
    check(err)

    var data map[string]interface{}

    err = json.Unmarshal(jsonBytes, &data)
    check(err)

    config.port = parsePort(data)
    config.time = parseTime(data)
}

func parsePort(data map[string]interface{}) int {
    if v,in := data["port"]; in {
        return int(v.(float64))
    }
    return DefaultPort
}

func parseTime(data map[string]interface{}) int {
    if v,in := data["time"]; in {
        return int(v.(float64))
    }
    return DefaultTime
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}
