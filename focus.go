package main

import (
    "flag"
    "fmt"
    "os"
    "jhurt/focus_proxy/strikeset"
)

func handleSSCommand(ssCommand *flag.FlagSet, ssAddSite *string, ssRemoveSite *string, ssStart *bool) {
    if os.Args[1] != "ss" {
        fmt.Printf("%v is not a valid subcommand\n", os.Args[1])
        os.Exit(1)
    }

    ssCommand.Parse(os.Args[2:])

    if *ssAddSite != "" {
        strikeset.AddToBlockedFile(*ssAddSite)
    }

    if *ssRemoveSite != "" {
    }
}

func main() {
    ssCommand := flag.NewFlagSet("ss", flag.ExitOnError)
    ssAddSite := ssCommand.String("addSite", "", "add site to strikeset")
    ssRemoveSite := ssCommand.String("removeSite", "", "remove a site from strikeset")
    ssStart := ssCommand.Bool("start", false, "start after editing strikeset")
    

    if len(os.Args) < 2 {
        // start proxy normally
    } else {
        handleSSCommand(ssCommand, ssAddSite, ssRemoveSite, ssStart)
    }

}
