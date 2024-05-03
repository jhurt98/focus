package portnumber

import (
    "bufio"
    "os"
    "strings"
)

var dir string = "./portnumber/portnumber.txt"
var Port string = ":"

func WritePortNumber(number string) {
    err := os.Truncate(dir, 0)
    check(err)

    file, err := os.Create(dir)
    check(err)

    _,err = file.WriteString(number)
    check(err) 
}

func ReadPortNumber() {
    file, err := os.Open(dir)
    check(err)

    defer file.Close()
    scanner := bufio.NewScanner(file)

    scanner.Scan()
    
    SetPortNumber(scanner.Text())
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func setPortNumber(number string) {
    var sb strings.Builder
    sb.WriteString(Port)
    sb.WriteString(number)

    Port = sb.String()
}
