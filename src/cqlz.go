package main

import (
    "os"
    "strconv"
	"errors"

	"github.com/gocql/gocql"
)

type Trans struct {
    BankId          string
    Id              gocql.UUID
    ChequeBankId    string
    ChequeId        gocql.UUID
    ChequeAmount    int
    ChequeDate      string
    ChequeImg       string
    FromAcc         string
    ToAcc           string
    Timestamp       int64
    Digsig          string
    Status          string
}

type Cheque struct {
    BankId          string
    Id              gocql.UUID
    Amount          int
    Date            string
    Img             string
}

var Session *gocql.Session

func initCStarSession() {
    cluster := gocql.NewCluster(cassandraConfig.host)
    cluster.Port = port(cassandraConfig.port)
	cluster.Keyspace = cassandraConfig.keyspace
	cluster.Consistency = consistancy(cassandraConfig.consistancy)

    s, err := cluster.CreateSession()
    if(err != nil) {
        println("Error cassandra session:", err.Error())
        os.Exit(1)
    }
    Session = s
}

func clearCStarSession() {
    Session.Close()
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

func createTrans(trans *Trans) error {
    q := `
        INSERT INTO trans (
            bank_id,
            id,
            cheque_bank_id,
            cheque_id,
            cheque_amount,
            cheque_date,
            cheque_img,
            from_acc,
            to_acc,
            timestamp,
            digsig,
            status
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    err := Session.Query(q, trans.BankId, trans.Id, trans.ChequeBankId, trans.ChequeId, trans.ChequeAmount, trans.ChequeDate, trans.ChequeImg, trans.FromAcc, trans.ToAcc, trans.Timestamp, trans.Digsig, trans.Status).Exec()
    if err != nil {
        println(err.Error())
    }

    return err
}

func createCheque(cheque *Cheque) error {
    q := `
        INSERT INTO cheques (
            bank_id,
            id,
            amount,
            date,
            img
        )
        VALUES (?, ?, ?, ?, ?)
    `
    err := Session.Query(q, cheque.BankId, cheque.Id, cheque.Amount, cheque.Date, cheque.Img).Exec()

    if err != nil {
        println(err.Error())
    }

    return err
}

func getCheque(bankId string, cId string)(*Cheque, error) {
    uuid, err := gocql.ParseUUID(cId)
    if err != nil {
        println(err.Error)
        return nil, err
    }

    m := map[string]interface{}{}
    q := `
        SELECT bank_id, id, amount, date, img
        FROM cheques
            WHERE bank_id = ?
            AND id = ?
        LIMIT 1
    `
    itr := Session.Query(q, bankId, uuid).Consistency(gocql.One).Iter()
    for itr.MapScan(m) {
        cheque := &Cheque{}
        cheque.BankId = m["bank_id"].(string)
        cheque.Id = m["id"].(gocql.UUID)
        cheque.Amount = m["amount"].(int)
        cheque.Date = m["date"].(string)
        cheque.Img = m["img"].(string)

        return cheque, nil
    }

    return nil, errors.New("Not found cheque")
}

func isDoubleSpend(from string, to string, cid string)bool {
    // parse cid and get uuid
    uuid, err := gocql.ParseUUID(cid)
    if err != nil {
        println(err.Error)
        return true
    }

    m := map[string]interface{}{}
    q := `
        SELECT id FROM trans
            WHERE from_acc=?
            AND cheque_id=?
        LIMIT 1
        ALLOW FILTERING
    `
    itr := Session.Query(q, from, uuid).Consistency(gocql.One).Iter()
    for itr.MapScan(m) {
        return true
    }

    q = `
        SELECT id FROM trans
            WHERE to_acc=?
            AND cheque_id=?
        LIMIT 1
        ALLOW FILTERING
    `
    itr = Session.Query(q, to, uuid).Consistency(gocql.One).Iter()
    for itr.MapScan(m) {
        return true
    }

    return false
}

func uuid() gocql.UUID {
    return gocql.TimeUUID()
}

func cuuid(cid string) (gocql.UUID, error) {
    return gocql.ParseUUID(cid)
}
