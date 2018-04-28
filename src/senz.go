package main

import (
	"bufio"
	"fmt"
	"net"
)

type Senzie struct {
	name   string
	out    chan string
	quit   chan bool
	tuk    chan string
	reader *bufio.Reader
	writer *bufio.Writer
	conn   *net.TCPConn
}

type Senz struct {
	Msg      string
	Uid      string
	Ztype    string
	Sender   string
	Receiver string
	Attr     map[string]string
	Digsig   string
}

// buffer size
const bufSize = 64 * 1024

func m() {
	// first init key pair
	setUpKeys()

	// init cassandra session
	initCStarSession()

	// address
	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.switchHost+":"+config.switchPort)
	if err != nil {
		fmt.Println("Error address:", err.Error())
		return
	}

	// tcp connect
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer conn.Close()

	// create senzie
	senzie := &Senzie{
		name:   config.senzieName,
		out:    make(chan string),
		quit:   make(chan bool),
		tuk:    make(chan string),
		reader: bufio.NewReaderSize(conn, bufSize),
		writer: bufio.NewWriterSize(conn, bufSize),
		conn:   conn,
	}
	registering(senzie)

	// close session
	clearCStarSession()
}

func registering(senzie *Senzie) {
	// send reg
	z := regSenz()
	senzie.writer.WriteString(z + ";")
	senzie.writer.Flush()

	// listen for reg status
	msg, err := senzie.reader.ReadString(';')
	if err != nil {
		fmt.Println("Error reading: ", err.Error())

		senzie.conn.Close()
		return
	}

	// parse senz
	// check reg status
	senz := parse(msg)
	if senz.Attr["status"] == "REG_DONE" || senz.Attr["status"] == "REG_ALR" {
		// reg done
		// start reading and writing
		go writing(senzie)
		reading(senzie)
	} else {
		// close and exit
		senzie.conn.Close()
		return
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
		if msg == "TAK;" {
			// when connect, we recive TAK
			continue READER
		} else if msg == "TIK;" {
			// send TIK
			senzie.tuk <- "TUK;"
			continue READER
		} else if msg == "TUK;" {
			continue READER
		}

		// handle msg via goroute
		go handling(senzie, msg)
	}
}

func writing(senzie *Senzie) {
WRITER:
	for {
		select {
		case <-senzie.quit:
			println("quiting/write -- ")
			break WRITER
		case senz := <-senzie.out:
			println("writing -- ")
			println(senz)
			// send
			senzie.writer.WriteString(senz + ";")
			senzie.writer.Flush()
		case tuk := <-senzie.tuk:
			senzie.writer.WriteString(tuk)
			senzie.writer.Flush()
		}
	}
}

func handling(senzie *Senzie, msg string) {
	// parse and handle
	senz := parse(msg)

	// TODO handle parse error

	// TODO verify signature

	// frist send AWA back
	senzie.out <- awaSenz(senz.Attr["uid"])

	if id, ok := senz.Attr["id"]; !ok {
		// this means new cheque
		// and new trans
		promize := senzToPromize(&senz)
		trans := senzToTrans(&senz, promize)
		trans.FromBank = senz.Attr["bnk"]
		trans.FromAccount = senz.Attr["acc"]
		trans.ToAccount = transConfig.account
		trans.Type = "TRANSFER"

		// call finacle to fund transfer
		err := doFundTrans(trans.FromAccount, trans.ToAccount, trans.PromizeAmount, transConfig.commission, senz.Attr["uid"])
		if err != nil {
			senzie.out <- statusSenz("ERROR", senz.Attr["uid"], senz.Sender)
			return
		}

		// create cheque
		// create trans
		createPromize(promize)
		createTrans(trans)

		// TODO handle create failures

		// send status back to fromAcc
		// forward cheque to toAcc
		senzie.out <- statusSenz("SUCCESS", senz.Attr["uid"], senz.Sender)
		senzie.out <- promizeSenz(promize, senz.Sender, senz.Attr["to"], uid())
	} else {
		// this mean already transfered cheque
		// check for double spend
		if isDoubleSpend(senz.Sender, id) {
			// send error status back
			senzie.out <- statusSenz("ERROR", senz.Attr["uid"], senz.Sender)
			return
		}

		// get cheque first
		promize, err := getPromize(config.SenzieName, id)
		if err != nil {
			senzie.out <- statusSenz("ERROR", senz.Attr["uid"], senz.Sender)
			return
		}

		// new trans
		trans := senzToTrans(&senz, promize)
		trans.FromAccount = transConfig.account
		trans.ToBank = senz.Attr["bnk"]
		trans.ToAccount = senz.Attr["acc"]
		trans.Type = "REDEEM"

		// call finacle to fund transfer
		err = doFundTrans(trans.FromAccount, trans.ToAccount, trans.PromizeAmount, transConfig.commission, id)
		if err != nil {
			senzie.out <- statusSenz("ERROR", senz.Attr["uid"], senz.Sender)
			return
		}

		// create trans
		createTrans(trans)

		// send success status back
		senzie.out <- statusSenz("SUCCESS", senz.Attr["uid"], senz.Sender)
	}
}
