package utils

import (
    "net"
    "net/http"
    "strings"
)

func GetIP(r *http.Request) string {
    ip := r.Header.Get("X-Forwarded-For")
    if ip != "" {
        return strings.Split(ip, ",")[0]
    }

    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return r.RemoteAddr
    }

    return host
}
