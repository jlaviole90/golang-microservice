package utils

import (
	"log"
	"os"
)

func Ptr[T any](v T) *T {
	return &v
}

func ReqEnvs(k string) string {
    v := os.Getenv(k)
    if v == "" {
        log.Fatalf("Missing required environment varaibles in connection. %s", k)
    }
    return v 
}
