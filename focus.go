package main

import (
    "flag"
)

func main() {
    blockSite := flag.Boolean("blockSite", false, "add a site to list of blocked sites")
    start := flag.Boolean("start", true, "force the proxy to start, used if you want to start after blocking a site")
    
}
