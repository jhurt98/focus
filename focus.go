package main

import (
    "flag"
    "fmt"
    "os"
    "jhurt/focus_proxy/internal/strikeset"
    "jhurt/focus_proxy/internal/config"
    "jhurt/focus_proxy/server"
)

func handleSSCommand(ssCommand *flag.FlagSet, ssAddSite *string, ssRemoveSite *string, ssStart *bool) {
    ssCommand.Parse(os.Args[2:])

    if *ssAddSite != "" {
        strikeset.AddToBlockedFile(*ssAddSite)
    }

    if *ssRemoveSite != "" {
        strikeset.RemoveFromBlockedFile(*ssRemoveSite)
    }

    if *ssStart { 
        proxy.StartProxy()
    }
}

func handleStartCommand(startCommand *flag.FlagSet, startTimer *int, startPort *int) {
    startCommand.Parse(os.Args[2:])
    config.Init()

    if *startTimer != -1 {
        config.SetTimeout(*startTimer)
    }
    
    if *startPort != 2002 {
        config.SetPortNumber(*startPort)
    }
    proxy.StartProxy()
}

func main() {
    ssCommand := flag.NewFlagSet("ss", flag.ExitOnError)
    ssAddSite := ssCommand.String("addSite", "", "add site to strikeset")
    ssRemoveSite := ssCommand.String("removeSite", "", "remove a site from strikeset")
    ssStart := ssCommand.Bool("start", false, "start after editing strikeset")
    

    startCommand := flag.NewFlagSet("start", flag.ExitOnError)
    startTimer := startCommand.Int("t", -1, "set timer in minutes")
    startPort := startCommand.Int("p", 2002, "set port number")

    if len(os.Args) < 2 {
        fmt.Println("expected 'ss' 'start' or 'config'")
        os.Exit(1)
    }

    switch os.Args[1] {

    case "ss":
        handleSSCommand(ssCommand, ssAddSite, ssRemoveSite, ssStart)
        return
    case "start":
        handleStartCommand(startCommand, startTimer, startPort)
        return
    default:
        fmt.Println("expected ss or start command")
        os.Exit(1)
    }

}
