package main

import (
    "os"
)

type Config struct {
	switchName      string
	switchHost      string
    switchPort      string
    senzieMode      string
    senzieName      string
    dotKeys         string
    idRsa           string
    idRsaPub        string
}

type CassandraConfig struct {
    host        string
    port        string
    keyspace    string
    consistancy string
}

type FinacleConfig struct {
    api             string
    lienAddAction   string
    lienModAction   string
}

var config = Config {
    switchName: getEnv("SWITCH_NAME", "senzswitch"),
    switchHost: getEnv("SWITCH_HOST", "www.rahasak.com"),
    switchPort: getEnv("SWITCH_PORT", "7070"),
    senzieMode: getEnv("SENZIE_MODE", "dev"),
    senzieName: getEnv("SENZIE_NAME", "sampath.chain"),
    dotKeys: getEnv("DOT_KEYS", ".keys"),
    idRsa: getEnv("ID_RSA", ".keys/id_rsa"),
    idRsaPub: getEnv("ID_RSA_PUB", ".keys/id_rsa.pub"),
}

var cassandraConfig = CassandraConfig {
    host: getEnv("CASSANDRA_HOST", "dev.localhost"),
    port: getEnv("CASSANDRA_PORT", "9042"),
    keyspace: getEnv("CASSANDRA_KEYSPACE", "cchain"),
    consistancy: getEnv("CASSANDRA_CONSISTANCY", "ALL"),
}

var finacleConfig = FinacleConfig {
    api: getEnv("FINACLE_API", "https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService"),
    lienAddAction: getEnv("LIEN_ADD_ACTION", "https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService"),
    lienModAction: getEnv("LIEN_MOD_ACTION", "https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService"),
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }

    return fallback
}
