package middleware

import (
    "context"
    "net/http"
    "strings"

    "dealls-test/utils"
)

func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
        claims, err := utils.ParseToken(tokenString)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Ambil request_id dari header (atau generate jika tidak ada)
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = utils.GenerateRequestID() // kamu bisa buat fungsi ini sendiri
        }

        // Simpan user ID, role, dan request_id ke context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "role", claims.Role)
        ctx = context.WithValue(ctx, "request_id", requestID)

        // Lanjut ke handler dengan context yang sudah disiapkan
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
