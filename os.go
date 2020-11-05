package main

import "os"

func GetEnvOfDefault(key, val string) string {
    if envvar := os.Getenv(key); len(envvar) > 0 {
        return envvar
    } else {
        return val
    }
}
