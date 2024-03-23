package proxy 

import (
    "fmt"
    "net"
    "net/http"
    "log"
    "context"
    "os/signal"
    "syscall"
    "os"
    "io"
    "jhurt/focus_proxy/strikeset"
)

var ss strikeset.Strikeset

func serviceRequest(w http.ResponseWriter, req *http.Request) {
    log.Printf("target host: %v\n", req.Host)
    log.Printf("request : %v\n", req)

    if req.RequestURI[0] == '/' {
        fmt.Fprintf(w, "proxy boi")
        return
    }

    if req.Method == http.MethodConnect { 
        handleCONNECTRequest(w, req)
        return
    } else {
        handleHTTP(w, req)
        return
    }
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
    if !ss.SiteIsAllowed(req.Host) {
        w.WriteHeader(http.StatusForbidden)
        fmt.Fprintf(w, "get to work")
        return
    }

    res, err := http.DefaultTransport.RoundTrip(req)
    if err != nil {
        panic(err)
    }

    defer res.Body.Close()
    copyHeader(w.Header(), res.Header)
    w.WriteHeader(res.StatusCode)
    io.Copy(w, res.Body)
    log.Printf("get: %v\n", res.Body)
}

func handleCONNECTRequest(w http.ResponseWriter, req *http.Request) {
    if !ss.SiteIsAllowed(req.Host) {
        w.WriteHeader(http.StatusForbidden)
        return 
    }

    targetConn, err := net.Dial("tcp", req.Host)
    if err != nil {
        log.Printf("woah woah something went wrong with dialing %v\n", req.Host)
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }

    w.WriteHeader(http.StatusOK)
    
    hijacker, ok := w.(http.Hijacker)
    if !ok {
        log.Fatal("response writer is not a hijacker")
    }

    clientCnxn, _, err := hijacker.Hijack()
    if err != nil { 
        log.Fatal("failed to hijack")
    }

    go tunnel(targetConn, clientCnxn)
    go tunnel(clientCnxn, targetConn)
}

func tunnel(dst io.WriteCloser, src io.ReadCloser) {
    io.Copy(dst, src)
    dst.Close()
    src.Close()
}

func copyHeader(dst, src http.Header) {
    for name, values := range src {
        for _, value := range values {
            dst.Add(name, value)
        }
    }
}

func startProxy() {
    ss = strikeset.Strikeset{}
    ss.AddToStrikeset()
    fmt.Printf("strikeset blocking: %v", ss.Domains)
    server := http.Server {
        Addr: ":2002",
        Handler: http.HandlerFunc(serviceRequest),
    }

    log.Printf("server listening on port: %v, pid: %v\n", server.Addr, os.Getpid())

    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

    defer cancel()

    go func() {
        err := server.ListenAndServe()
        if err != nil && err != http.ErrServerClosed{
            panic(err)
        }
    }()

    <-ctx.Done()
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("server shutdown fail")
    }
}
