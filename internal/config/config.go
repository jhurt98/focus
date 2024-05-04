package config
import (
    "os"
    "time"
    "encoding/json"
    "strconv"
)

type Config struct {
    Port string 
    Time time.Duration 
}

var DefaultPort = 20002
var DefaultTime time.Duration  = -1

var configJsonFile = "./internal/config/config.json"

var MainConfig *Config

func Init() {
    MainConfig = buildConfig() 
}

func buildConfig() *Config {
    config := &Config{}
    config.setValuesWithJson()
    return config
}

func (config *Config) setValuesWithJson() {
    data := readConfigJson()
    config.Port = parsePort(data)
    config.Time = parseTime(data)
}

func readConfigJson() map[string]interface{} {
    jsonBytes, err := os.ReadFile(configJsonFile)
    check(err)

    var data map[string]interface{}

    err = json.Unmarshal(jsonBytes, &data)
    check(err)

    return data
}

func parsePort(data map[string]interface{}) string {
    portNumber := DefaultPort
    if v,in := data["port"]; in {
        vF, ok := v.(float64)
        if ok {
            portNumber = int(vF)
        }
    }
    return buildPortString(portNumber) 
}

func parseTime(data map[string]interface{}) time.Duration {
    if v,in := data["time"]; in {
        vF, ok := v.(float64)
        if ok {
            return time.Duration(int(vF))
        }
    }
    return DefaultTime
}

func SetPortNumber(number int) {
    if MainConfig == nil {
        panic("config has not been initialized appropriately")
    }

    MainConfig.Port = buildPortString(number)
}

func SetTimeout(timeout int) {
    if MainConfig == nil {
        panic("config has not been initialized appropriately")
    }
    MainConfig.Time = time.Duration(timeout)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func buildPortString(portNumber int) string {
    Port := ":"
    return Port + strconv.Itoa(portNumber)
}
