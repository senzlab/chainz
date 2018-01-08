package main

import (
    "os"
)

type Config struct {
	switchName  string
	switchHost  string
    switchPort  string
    senzieMode  string
    senzieName  string
    dotKeys     string
    idRsa       string
    idRsaPub    string
}

var config = Config {
    switchName: getEnv("SWITCH_NAME", "senzswitch"),
    switchHost: getEnv("SWITCH_HOST", "www.rahasak.com"),
    switchPort: getEnv("SWITCH_PORT", "7070"),
    senzieMode: getEnv("SENZIE_MODE", "dev"),
    senzieName: getEnv("SENZIE_NAME", "sampath"),
    dotKeys: getEnv("DOT_KEYS", ".keys"),
    idRsa: getEnv("ID_RSA", ".keys/id_rsa"),
    idRsaPub: getEnv("ID_RSA_PUB", ".keys/id_rsa.pub"),
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }

    return fallback
}
