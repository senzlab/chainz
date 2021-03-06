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
	api               string
	action            string
	commission        string
	commissionAccount string
	account           string
	bankCode          string
	ceftCommission    string
}

type SmzConfig struct {
	appcode  string
	password string
	api      string
}

var config = Config{
	switchName: getEnv("SWITCH_NAME", "senzswitch"),
	switchHost: getEnv("SWITCH_HOST", "www.rahasak.com"),
	switchPort: getEnv("SWITCH_PORT", "7070"),
	senzieMode: getEnv("SENZIE_MODE", "dev"),
	senzieName: getEnv("SENZIE_NAME", "sampath"),
	dotKeys:    getEnv("DOT_KEYS", ".keys"),
	idRsa:      getEnv("ID_RSA", ".keys/id_rsa"),
	idRsaPub:   getEnv("ID_RSA_PUB", ".keys/id_rsa.pub"),
}

var cassandraConfig = CassandraConfig{
	host:        getEnv("CASSANDRA_HOST", "dev.localhost"),
	port:        getEnv("CASSANDRA_PORT", "9042"),
	keyspace:    getEnv("CASSANDRA_KEYSPACE", "zchain"),
	consistancy: getEnv("CASSANDRA_CONSISTANCY", "ALL"),
}

var transConfig = TransConfig{
	api:               getEnv("TRANS_API", "http://192.168.125.93:7800/sd/iib/IIBFinacleIntegration"),
	action:            getEnv("TRANS_ACTION", "http://192.168.125.93:7800/sd/iib/iibfinacleintegration"),
	commission:        getEnv("TRANS_COMMISSION", "20"),
	commissionAccount: getEnv("COMMISSION_ACCOUNT", "900108020041"),
	account:           getEnv("TRANS_ACCOUNT", "900100000801"),
	bankCode:          getEnv("TRANS_BANK_CODE", "7278"),
	ceftCommission:    getEnv("CEFT_COMMISSION", "50"),
}

var smzConfig = SmzConfig{
	appcode:  getEnv("SMZ_APPCODE", "IGIFT"),
	password: getEnv("SMZ_PASSWORD", "476d5ffbe2bb46196cc9a58ef8aec1bd"),
	api:      getEnv("SMZ_API", "http://192.125.125.143:7005/SMSSolution_V1.02/SampathSMS/special"),
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
