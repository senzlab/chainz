package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
)

type Senzie struct {
    name        string
	out         chan Senz
    quit        chan bool
    tik         chan string
    reader      *bufio.Reader
    writer      *bufio.Writer
    conn        *net.TCPConn
}

type Senz struct {
    Msg         string
    Uid          string
    Ztype       string
    Sender      string
    Receiver    string
    Attr        map[string]string
    Digsig      string
}

func main() {
    // first init key pair
    setUpKeys()

    // init cassandra session
    initCStarSession()

    // init trans
    trans := &Trans{}
    trans.BankId = "sampath"
    trans.ChequeBankId = "sampath"
    trans.ChequeAmount = 1000
    trans.ChequeDate = "12/01/2018"
    trans.FromAcc = "111111"
    trans.ToAcc = "222222"
    trans.Digsig = "DIGISIG"
    trans.Status = "PENDING"
    createTrans(trans)
    println(isDoubleSpend("111111", "222222", "768be756-f51f-11e7-a0f7-4c327597ce77"))

    // address
    //tcpAddr, err := net.ResolveTCPAddr("tcp4", config.switchHost + ":" + config.switchPort)
    //if err != nil {
    //    fmt.Println("Error address:", err.Error())
    //    os.Exit(1)
    //}

    // tcp connect
    //conn, err := net.DialTCP("tcp", nil, tcpAddr)
    //if err != nil {
    //    fmt.Println("Error listening:", err.Error())
    //    os.Exit(1)
    //}

    // close on app closes
    //defer conn.Close()

    //fmt.Println("connected to switch")

    // create senzie
    //senzie := &Senzie {
    //    name: config.senzieName,
    //    out: make(chan Senz),
    //    quit: make(chan bool),
    //    tik: make(chan string),
    //    reader: bufio.NewReader(conn),
    //    writer: bufio.NewWriter(conn),
    //    conn: conn,
    //}
    //registering(senzie)

    // close session
    clearCStarSession()
}

func registering(senzie *Senzie) {
    // send reg
    uid := uid()
    pubkey := getIdRsaPubStr() 
    z := "SHARE #pubkey " + pubkey +
                " #uid " + uid +
                " @" + config.switchName +
                " ^" + config.senzieName +
                " digisig"
    senzie.writer.WriteString(z + ";")
    senzie.writer.Flush()

    // listen for reg status
    msg, err := senzie.reader.ReadString(';')
    if err != nil {
        fmt.Println("Error reading: ", err.Error())

        senzie.conn.Close()
        os.Exit(1)
    }

    // parse senz
    // check reg status
    senz := parse(msg)
    if(senz.Attr["status"] == "REG_DONE" || senz.Attr["status"] == "REG_ALR") {
        println("reg done...")
        // start reading and writing
        go writing(senzie)
        reading(senzie)
    } else {
        // close and exit
        senzie.conn.Close()
        os.Exit(1)
    }
}

func reading(senzie *Senzie) {
    READER:
    for {
        // read data
        msg, err := senzie.reader.ReadString(';')
        if err != nil {
            fmt.Println("Error reading: ", err.Error())

            senzie.quit <- true
            break READER
        }

        // not handle TAK, TIK, TUK
        if (msg == "TAK;") {
            // when connect, we recive TAK
            continue READER
        } else if(msg == "TIK;") {
            // send TIK
            senzie.tik <- "TUK;"
            continue READER
        } else if(msg == "TUK;") {
            continue READER
        }

        println("---- " + msg)

        // parse and handle
        senz := parse(msg)
        handling(senzie, &senz)
    }
}

func writing(senzie *Senzie)  {
    // write
    WRITER:
    for {
        select {
        case <- senzie.quit:
            println("quiting/write -- ")
            break WRITER
        case senz := <-senzie.out:
            println("writing -- ")
            println(senz.Msg)
            // TODO sign and send
            senzie.writer.WriteString(senz.Msg + ";")
            senzie.writer.Flush()
        case tik := <- senzie.tik:
            println("ticking -- " )
            senzie.writer.WriteString(tik)
            senzie.writer.Flush()
        }
    }
}

func handling(senzie *Senzie, senz *Senz) {
    // frist send AWA back
    uid := senz.Attr["uid"]
    z := "AWA #uid " + uid + 
              " @" + config.switchName +
              " ^" + config.senzieName +
              " digisig"
    sz := Senz{}
    sz.Uid = uid
    sz.Msg = z
    sz.Receiver = config.switchName
    sz.Sender = config.senzieName
    senzie.out <- sz

    if(senz.Ztype == "SHARE") {
        // we only handle share cheques
        // get cheque attributes
        bId := config.senzieName
        cBnkId := senz.Attr["cbank"]
        cId := senz.Attr["cid"]
        //cAmnt := senz.Attr["camnt"]
        cDate := senz.Attr["cdate"]
        cImg := senz.Attr["cimg"]
        toAcc := senz.Attr["to"]
        fromAcc := senz.Sender
        digsig := senz.Digsig

        if (len(cId) == 0) {
            // this means new cheque
            // create cheque
            cheque := &Cheque{}
            cheque.BankId = bId;
            cheque.Id = uuid()
            cheque.Amount = 1000
            cheque.Date = cDate
            cheque.Img = cImg
            createCheque(cheque)

            // create trans
            trans := &Trans{}
            trans.BankId = bId
            trans.Id = uuid()
            trans.ChequeBankId = cBnkId
            trans.ChequeId = cheque.Id
            trans.ChequeAmount = 1000
            trans.ChequeDate = cDate
            trans.ChequeImg = cImg
            trans.FromAcc = fromAcc
            trans.ToAcc = toAcc
            trans.Digsig = digsig
            trans.Status = "TRANSFER"
            createTrans(trans)

            // TODO forward cheque to toAcc
            // TODO send status back to fromAcc
        } else {
            // this mean already transfered cheque
            // check for double spend
            if(isDoubleSpend(fromAcc, toAcc, cId)) {
                // TODO send error status back
            } else {
                // TODO create trans 
            }
        }
    }
}
