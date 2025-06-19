package utils

import (
    "github.com/google/uuid"
)

func GenerateRequestID() string {
    return "req-" + uuid.New().String()
}
