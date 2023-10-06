package main

import (
    "log"
    "strings"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
)

func makeHandler(upstream string) func (w http.ResponseWriter, r *http.Request){
    url, _ := url.Parse(upstream)
    proxy := httputil.NewSingleHostReverseProxy(url)

    return func (w http.ResponseWriter, r *http.Request) {
        if prior, ok := r.Header["X-Forwarded-For"]; ok {
            r.Header.Set("X-Real-IP", strings.Join(prior, ", "))
        }
        r.Header.Set("X-Proxy-Version", "gae-reverse-proxy/0.0.1")
        r.URL.Host = url.Host
        r.URL.Scheme = url.Scheme
        r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
        proxy.ServeHTTP(w, r)
    }
}

func main() {
    port := os.Getenv("PORT")
    upstream := os.Getenv("UPSTREAM")
    http.HandleFunc("/", makeHandler(upstream))
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
