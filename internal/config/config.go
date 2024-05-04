package config
import (
    "fmt"
    "os"
    "encoding/json"
    "strconv"
)

type Config struct {
    port string 
    time int
}

var DefaultPort = 20002
var DefaultTime = -1

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
    config.port = parsePort(data)
    config.time = parseTime(data)
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

func parseTime(data map[string]interface{}) int {
    if v,in := data["time"]; in {
        vF, ok := v.(float64)
        if ok {
            return int(vF)
        }
    }
    return DefaultTime
}

func WritePortNumber(number float64) {
    jsonBytes, err := os.ReadFile(configJsonFile)
    check(err)

    var data map[string]interface{}

    err = json.Unmarshal(jsonBytes, &data)
    check(err)

    data["port"] = number
    res,err := json.Marshal(data)
    check(err)

    fmt.Printf("resulted marshalled json %s\n", res)
}

func WriteTime() {
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
