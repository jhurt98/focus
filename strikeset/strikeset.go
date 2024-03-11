package strikeset

import ( 
    "os"
    "bufio"
    "regexp"
)

var dir string = "./strikeset/blocked.txt"

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

func (ss *Strikeset) AddToStrikeset() {
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
