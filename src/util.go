package main

import (
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
	return "1.24"
}

func senzToPromize(senz *Senz) *Promize {
	promize := &Promize{}
	promize.Bank = config.senzieName
	promize.Id = uuid()
	promize.Amount = senz.Attr["amnt"]
	promize.Blob = senz.Attr["blob"]
	promize.OriginZaddress = senz.Sender
	promize.OriginBank = senz.Attr["bnk"]
	promize.OriginAccount = senz.Attr["acc"]

	return promize
}

func senzToTrans(senz *Senz, promize *Promize) *Trans {
	trans := &Trans{}
	trans.Bank = config.senzieName
	trans.Id = uuid()
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
	user.Bank = senz.Attr["bank"]
	user.Account = senz.Attr["account"]
	user.PublicKey = senz.Attr["pubkey"]

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
