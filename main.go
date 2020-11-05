package main

import (
    "log"
    "net/http"
)

const (
    RedisHostKey = "TODO_REDIS_HOST"
    RedisPortKey = "TODO_REDIS_PORT"
    RedisPassKey = "TODO_REDIS_PASSWORD"
)

func main() {
    http.Handle("/", handleIndex())
    http.Handle("/add", handleAdd())
    http.Handle("/done", handleDone())
    http.Handle("/del", handleDel())

    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Szerver indítása sikertelen. %v\n", err)
    }
}
