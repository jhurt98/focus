package strikeset

import ( 
    "fmt"
    "os"
    "bufio"
    "regexp"
)

var dir string = "./internal/strikeset/blocked.txt"

type Strikeset struct {
    Domains []*regexp.Regexp
}

func AddToBlockedFile(site string) {
    file, err := os.OpenFile(dir, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
    check(err)

    _, err = file.WriteString(site)
    check(err)

    _, err = file.WriteString("\n")
    check(err)

    check(file.Close())
}

func RemoveFromBlockedFile(site string) {
     count := 0
     var lines []string

     file, err := os.Open(dir)
     check(err)

     defer file.Close()
     scanner := bufio.NewScanner(file)

     var line string
     for scanner.Scan() {
         line = scanner.Text()
         if line == site {
             count++
             continue
         } else {
             lines = append(lines, line)
         }
     }
     check(scanner.Err())

     if count == 0 {
         fmt.Printf("blocked.txt did not contain %v\n", site)
         return
     }

     err = os.Truncate(dir, 0)
     check(err)
     file, err = os.Create(dir)
     check(err)
     for _, line := range(lines) {
         _, err = file.WriteString(line)
         check(err)
         _, err = file.WriteString("\n")
         check(err)
     }
}

func (ss *Strikeset) LoadStrikeset() {
    file, err := os.Open(dir)
    check(err)

    defer file.Close()
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        ss.Domains = append(ss.Domains, getRegExp(scanner.Text()))
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
}

func (ss *Strikeset) SiteIsAllowed(site string) bool {
    for _, r := range(ss.Domains) {
        if r.MatchString(site) {
            return false 
        }
    }
    return true 
}

func getRegExp(word string) *regexp.Regexp{
    r, _ := regexp.Compile(word)
    return r;
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}
