package proxy 

import (
    "fmt"
    "net/http"
    "log"
    "context"
    "os/signal"
    "syscall"
    "os"
    "io"
    "net"
    "time"
    "jhurt/focus_proxy/internal/strikeset"
    "jhurt/focus_proxy/internal/config"
)

var ss strikeset.Strikeset

func serviceRequest(w http.ResponseWriter, req *http.Request) {
    if req.RequestURI[0] == '/' {
        fmt.Fprintf(w, "proxy")
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
}

func handleCONNECTRequest(w http.ResponseWriter, req *http.Request) {
    if !ss.SiteIsAllowed(req.Host) {
        handleBlockedHttps(w, req)
        return 
    }
    handleSafeHttps(w, req)

}

func handleBlockedHttps(w http.ResponseWriter, req *http.Request) {
    caCertFile, caKeyFile := GetCertAndKeyFiles()
    mitm := CreateMitmProxy(caCertFile, caKeyFile) 
    mitm.HandleTLS(w, req)
}

func handleSafeHttps(w http.ResponseWriter, req *http.Request) {
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

func StartProxy() {
    ss = strikeset.Strikeset{}
    ss.LoadStrikeset()

    port    := config.MainConfig.Port
    timeOut := config.MainConfig.Time

    server := http.Server {
        Addr: port,
        Handler: http.HandlerFunc(serviceRequest),
    }

    log.Printf("server listening on port: %v, pid: %v\n", server.Addr, os.Getpid())

    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

    defer cancel()

    timer := time.NewTimer(timeOut*time.Minute)
    if timeOut == -1 {
        if !timer.Stop() {
            <-timer.C
        }
    }

    go func() {
        err := server.ListenAndServe()
        if err != nil && err != http.ErrServerClosed{
            panic(err)
        }
    }()

    select {
    case <-timer.C:
        log.Println("timer went off")
        cancel()
    case <-ctx.Done():
        cancel()
    }

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("server shutdown fail", err)
    }
}

