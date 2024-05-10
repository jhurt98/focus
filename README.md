Focus CLI Proxy

CLI Tool that starts a proxy to block connections to specified sites

Commands:

"focus start": starts a proxy server at default port number until terminated by C-c

--> optional flags:

    --> "-t" sets a timeout in minutes for the current server session
    --> "-p" sets a custom port number for the current server session

"focus ss": command to edit the strikeset file for blocked sites/domains

--> flags:

    --> "addSite" adds a site to the file
    --> "removeSite" removes a site to the file if it exists

Depends on existing certificate authority and key files generated by the mkcert tool
- This is only used to spin a TLS server that accepts our local CA on https requests of blocked sites so our proxy server can interecept and return a custom response
- allowed sites do not create a TLS server and instead creates a blind tunnel between the client and server

Only tested on firefox by configuring the proxy in network settings
