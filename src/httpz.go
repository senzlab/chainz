package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// first init key pair
	setUpKeys()

	// init cassandra session
	initCStarSession()

	http.HandleFunc("/promize", promize)
	address := ":7070"
	println("starting server on address" + address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		println(err.Error)
		os.Exit(1)
	}
}

func promize(w http.ResponseWriter, r *http.Request) {
	// read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	msg := string(b)
	println(msg)

	// handle senz
	senz := parse(msg)
	if senz.Ztype == "SHARE" {
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
			err := doFundTrans(trans.FromAccount, trans.ToAccount, trans.PromizeAmount)
			if err != nil {
				resp := statusSenz("ERROR", senz.Attr["uid"], id, "cbid", senz.Sender)
				http.Error(w, resp, 400)
				return
			}

			// create cheque
			// create trans
			createPromize(promize)
			createTrans(trans)

			// TODO handle create failures

			// forward cheque to toAcc
			resp := promizeSenz(promize, senz.Sender, senz.Attr["to"], uid())
			w.WriteHeader(http.StatusCreated)
			io.WriteString(w, resp)
		} else {
			// this mean already transfered cheque
			// check for double spend
			if isDoubleSpend(senz.Sender, id) {
				// send error status back
				resp := statusSenz("ERROR", senz.Attr["uid"], id, "cbid", senz.Sender)
				http.Error(w, resp, 400)
				return
			} else {
				// get cheque first
				promize, err := getPromize(senz.Attr["bnk"], id)
				if err != nil {
					resp := statusSenz("ERROR", senz.Attr["uid"], id, "cbid", senz.Sender)
					http.Error(w, resp, 404)
					return
				} else {
					// new trans
					trans := senzToTrans(&senz, promize)
					trans.FromAccount = transConfig.account
					trans.ToBank = senz.Attr["bnk"]
					trans.ToAccount = senz.Attr["acc"]
					trans.Type = "REDEEM"

					// call finacle to fund transfer
					err := doFundTrans(trans.FromAccount, trans.ToAccount, trans.PromizeAmount)
					if err != nil {
						resp := statusSenz("ERROR", senz.Attr["uid"], id, "cbid", senz.Sender)
						http.Error(w, resp, 403)
						return
					}

					// create trans
					createTrans(trans)

					// send success status back
					resp := statusSenz("SUCCESS", senz.Attr["uid"], id, "cbid", senz.Sender)
					io.WriteString(w, resp)
				}
			}
		}
	} else {
		http.Error(w, "", 405)
		return
	}
}
