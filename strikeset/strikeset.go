package strikeset

import ( 
    "os"
    "bufio"
    "regexp"
)

type Strikeset struct {
    Domains []*regexp.Regexp
}

func (ss *Strikeset) AddToStrikeset(dir string) {
    file, err := os.Open(dir)
    if err != nil {
        panic(err)
    }
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
