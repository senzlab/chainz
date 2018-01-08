package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
    "strings"
    "strconv"
    "time"
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

    // address
    tcpAddr, err := net.ResolveTCPAddr("tcp4", config.switchHost + ":" + config.switchPort)
    if err != nil {
        fmt.Println("Error address:", err.Error())
        os.Exit(1)
    }

    // tcp connect
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    // close on app closes
    defer conn.Close()

    fmt.Println("connected to switch")

    // create senzie
    senzie := &Senzie {
        name: config.senzieName,
        out: make(chan Senz),
        quit: make(chan bool),
        tik: make(chan string),
        reader: bufio.NewReader(conn),
        writer: bufio.NewWriter(conn),
        conn: conn,
    }
    registering(senzie)
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
            
    } else if(senz.Ztype == "PUT") {
        
    }
}


func parse(msg string)Senz {
    fMsg := formatToParse(msg)
    tokens := strings.Split(fMsg, " ")
    senz := Senz {}
    senz.Msg = fMsg
    senz.Attr = map[string]string{}

    for i := 0; i < len(tokens); i++ {
        if(i == 0) {
            senz.Ztype = tokens[i]
        } else if(i == len(tokens) - 1) {
            // signature at the end
            senz.Digsig = tokens[i]
        } else if(strings.HasPrefix(tokens[i], "@")) {
            // receiver @eranga
            senz.Receiver = tokens[i][1:]
        } else if(strings.HasPrefix(tokens[i], "^")) {
            // sender ^lakmal
            senz.Sender = tokens[i][1:]
        } else if(strings.HasPrefix(tokens[i], "$")) {
            // $key er2232
            key := tokens[i][1:]
            val := tokens[i + 1]
            senz.Attr[key] = val
            i ++
        } else if(strings.HasPrefix(tokens[i], "#")) {
            key := tokens[i][1:]
            nxt := tokens[i + 1]

            if(strings.HasPrefix(nxt, "#") || strings.HasPrefix(nxt, "$") ||
                                                strings.HasPrefix(nxt, "@")) {
                // #lat #lon
                // #lat @eranga
                // #lat $key 32eewew
                senz.Attr[key] = ""
            } else {
                // #lat 3.2323 #lon 5.3434
                senz.Attr[key] = nxt
                i ++
            }
        }
    }

    // set uid as the senz id
    senz.Uid = senz.Attr["uid"]

    return senz
}

func uid()string {
    t := time.Now().UnixNano() / int64(time.Millisecond)
    return config.senzieName + strconv.FormatInt(t, 10)
}
