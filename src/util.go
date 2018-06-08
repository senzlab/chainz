package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func parse(msg string) Senz {
	fMsg := formatToParse(msg)
	tokens := strings.Split(fMsg, " ")
	senz := Senz{}
	senz.Msg = fMsg
	senz.Attr = map[string]string{}

	for i := 0; i < len(tokens); i++ {
		if i == 0 {
			senz.Ztype = tokens[i]
		} else if i == len(tokens)-1 {
			// signature at the end
			senz.Digsig = tokens[i]
		} else if strings.HasPrefix(tokens[i], "@") {
			// receiver @eranga
			senz.Receiver = tokens[i][1:]
		} else if strings.HasPrefix(tokens[i], "^") {
			// sender ^lakmal
			senz.Sender = tokens[i][1:]
		} else if strings.HasPrefix(tokens[i], "$") {
			// $key er2232
			key := tokens[i][1:]
			val := tokens[i+1]
			senz.Attr[key] = val
			i++
		} else if strings.HasPrefix(tokens[i], "#") {
			key := tokens[i][1:]
			nxt := tokens[i+1]

			if strings.HasPrefix(nxt, "#") || strings.HasPrefix(nxt, "$") ||
				strings.HasPrefix(nxt, "@") {
				// #lat #lon
				// #lat @eranga
				// #lat $key 32eewew
				senz.Attr[key] = ""
			} else {
				// #lat 3.2323 #lon 5.3434
				senz.Attr[key] = nxt
				i++
			}
		}
	}

	// set uid as the senz id
	senz.Uid = senz.Attr["uid"]

	return senz
}

func formatToParse(msg string) string {
	replacer := strings.NewReplacer(";", "", "\n", "", "\r", "")
	return strings.TrimSpace(replacer.Replace(msg))
}

func formatToSign(msg string) string {
	replacer := strings.NewReplacer(";", "", "\n", "", "\r", "", " ", "")
	return strings.TrimSpace(replacer.Replace(msg))
}

func uid() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return config.senzieName + strconv.FormatInt(t, 10)
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func randomSalt() string {
	// gengerate random salt between 0.10 - 2.00
	rand.Seed(time.Now().Unix())
	x := rand.Intn(2)
	y := rand.Intn(9-1) + 1
	z := rand.Intn(9)
	return fmt.Sprintf("%d.%d%d", x, y, z)
}

func commission(amount string) string {
	amnt, err := strconv.Atoi(amount)
	if err != nil {
		return ""
	}

	if amnt > 5000 {
		return "50"
	} else {
		return "20"
	}
}

func zode() string {
	l := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = l[rand.Intn(len(l))]
	}

	return string(b)
}

func senzToPromize(senz *Senz) *Promize {
	promize := &Promize{}
	promize.Bank = config.senzieName
	id, _ := cuuid(senz.Attr["id"])
	promize.Id = id
	promize.Amount = senz.Attr["amnt"]
	promize.Blob = senz.Attr["blob"]
	promize.OriginZaddress = senz.Sender
	promize.OriginBank = config.senzieName
	promize.OriginAccount = senz.Attr["acc"]
	promize.Timestamp = timestamp()

	return promize
}

func senzToTrans(senz *Senz, promize *Promize) *Trans {
	trans := &Trans{}
	trans.Bank = config.senzieName
	trans.Id = uuid()
	trans.Timestamp = timestamp()
	trans.PromizeBank = promize.Bank
	trans.PromizeId = promize.Id
	trans.PromizeAmount = promize.Amount
	trans.PromizeBlob = promize.Blob
	trans.FromZaddress = senz.Sender
	trans.ToZaddress = senz.Attr["to"]
	trans.Digsig = senz.Digsig

	return trans
}

func senzToUser(senz *Senz) *User {
	user := &User{}
	user.Zaddress = senz.Sender
	user.Bank = config.senzieName
	user.PublicKey = senz.Attr["pubkey"]
	user.Verified = false
	user.Active = true
	user.Timestamp = timestamp()

	return user
}

func regSenz() string {
	z := "SHARE #pubkey " + getIdRsaPubStr() +
		" #uid " + uid() +
		" @" + config.switchName +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func awaSenz(uid string) string {
	z := "AWA #uid " + uid +
		" @" + config.switchName +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func statusSenz(status string, uid string, to string) string {
	z := "DATA #status " + status +
		" #uid " + uid +
		" @" + to +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func zodeSenz(status string, zode string, uid string, to string) string {
	z := "DATA #status " + status +
		" #zode " + zode +
		" #uid " + uid +
		" @" + to +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func promizeSenz(promize *Promize, from string, to string, uid string) string {
	z := "SHARE #bnk " + promize.Bank +
		" #id " + promize.Id.String() +
		" #amnt " + promize.Amount +
		" #blob " + promize.Blob +
		" #from " + from +
		" #uid " + uid +
		" @" + to +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}
