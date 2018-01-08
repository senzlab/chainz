package main

import (
    "fmt"
    "os"
    "strconv"
	"github.com/gocql/gocql"
)

var Session *gocql.Session

func initCStar() {
    cluster := gocql.NewCluster(cassandraConfig.host)
    cluster.Port = port(cassandraConfig.port)
	cluster.Keyspace = cassandraConfig.keyspace
	cluster.Consistency = consistancy(cassandraConfig.consistancy)

    Session, err := cluster.CreateSession()
    if err != nil {
        fmt.Println("Error connecting to cassandra:", err.Error())
        os.Exit(1)
    }
	defer Session.Close()
}

func port(p string) int {
    i, err := strconv.Atoi(p)
    if err != nil {
        return 9042
    }

    return i
}

func consistancy(c string) gocql.Consistency {
    gc, err := gocql.MustParseConsistency(c)
    if err != nil {
        return gocql.All
    }

    return gc
}


