package middleware

import (
    "net/http"
    "strings"
    "dealls-test/utils"
)

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
        claims, err := utils.ParseToken(tokenString)
        if err != nil || claims.Role != "admin" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
 
        next.ServeHTTP(w, r)
    }
}
