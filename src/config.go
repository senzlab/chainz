package main

import (
	"os"
)

type Config struct {
	switchName string
	switchHost string
	switchPort string
	senzieMode string
	senzieName string
	dotKeys    string
	idRsa      string
	idRsaPub   string
}

type CassandraConfig struct {
	host        string
	port        string
	keyspace    string
	consistancy string
}

type TransConfig struct {
	api        string
	action     string
	commission string
	account    string
}

var config = Config{
	switchName: getEnv("SWITCH_NAME", "senzswitch"),
	switchHost: getEnv("SWITCH_HOST", "www.rahasak.com"),
	switchPort: getEnv("SWITCH_PORT", "7070"),
	senzieMode: getEnv("SENZIE_MODE", "dev"),
	senzieName: getEnv("SENZIE_NAME", "sampath.chain"),
	dotKeys:    getEnv("DOT_KEYS", ".keys"),
	idRsa:      getEnv("ID_RSA", ".keys/id_rsa"),
	idRsaPub:   getEnv("ID_RSA_PUB", ".keys/id_rsa.pub"),
}

var cassandraConfig = CassandraConfig{
	host:        getEnv("CASSANDRA_HOST", "dev.localhost"),
	port:        getEnv("CASSANDRA_PORT", "9042"),
	keyspace:    getEnv("CASSANDRA_KEYSPACE", "cchain"),
	consistancy: getEnv("CASSANDRA_CONSISTANCY", "ALL"),
}

var transConfig = TransConfig{
	api:        getEnv("TRANS_API", "https://uatweb.sampath.lk/SVRClientWeb/services/MobileServicesMain"),
	action:     getEnv("TRANS_ACTION", "http://mobileServices.web.app.sampath.org/MobileServicesMain/initiateTransferRequest"),
	commission: getEnv("TRANS_COMMISSION", "20"),
	account:    getEnv("TRANS_ACCOUNT", "231233223344"),
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
