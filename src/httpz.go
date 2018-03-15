package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type SenzMsg struct {
	Uid string
	Msg string
}

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

	println(string(b))

	// unmarshel json
	var senzMsg SenzMsg
	err = json.Unmarshal(b, &senzMsg)
	if err != nil {
		resp := statusSenz("ERROR", "uid", "id", "cbid", "sender")
		http.Error(w, resp, 400)
		return
	}

	// handle senz
	senz := parse(senzMsg.Msg)
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
				// marshel and return error
				senzMsg := SenzMsg{
					Uid: senz.Attr["uid"],
					Msg: statusSenz("ERROR", senz.Attr["uid"], id, "cbid", senz.Sender),
				}
				j, _ := json.Marshal(senzMsg)
				http.Error(w, string(j), 400)
			}

			// create cheque
			// create trans
			createPromize(promize)
			createTrans(trans)

			// TODO handle create failures

			// msg to sender
			fMsg := SenzMsg{
				Uid: senz.Attr["uid"],
				Msg: statusSenz("SUCCESS", senz.Attr["uid"], "id", "sampath", senz.Sender),
			}

			// message for #to
			tMsg := SenzMsg{
				Uid: senz.Attr["uid"],
				Msg: promizeSenz(promize, senz.Sender, senz.Attr["to"], uid()),
			}

			// marshel and return response(forward chque to toAcc senz)
			var zmsgs []SenzMsg
			zmsgs = append(zmsgs, fMsg)
			zmsgs = append(zmsgs, tMsg)
			j, _ := json.Marshal(zmsgs)
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, string(j))
		} else {
			// this mean already transfered cheque
			// check for double spend
			if isDoubleSpend(senz.Sender, id) {
				// marshel and return error
				senzMsg := SenzMsg{
					Uid: senz.Attr["uid"],
					Msg: statusSenz("ERROR", senz.Attr["uid"], id, "cbid", senz.Sender),
				}
				j, _ := json.Marshal(senzMsg)
				http.Error(w, string(j), 400)
				return
			} else {
				// get cheque first
				promize, err := getPromize(senz.Attr["bnk"], id)
				if err != nil {
					// marshel and return error
					senzMsg := SenzMsg{
						Uid: senz.Attr["uid"],
						Msg: statusSenz("ERROR", senz.Attr["uid"], id, "cbid", senz.Sender),
					}
					j, _ := json.Marshal(senzMsg)
					http.Error(w, string(j), 404)
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
						// marshel and return error
						senzMsg := SenzMsg{
							Uid: senz.Attr["uid"],
							Msg: statusSenz("ERROR", senz.Attr["uid"], id, "cbid", senz.Sender),
						}
						j, _ := json.Marshal(senzMsg)
						http.Error(w, string(j), 403)
						return
					}

					// create trans
					createTrans(trans)

					// status to sender
					zmsg := SenzMsg{
						Uid: senz.Attr["uid"],
						Msg: statusSenz("SUCCESS", senz.Attr["uid"], id, "cbid", senz.Sender),
					}

					// marshel and return response(success response to sender)
					var zmsgs []SenzMsg
					zmsgs = append(zmsgs, zmsg)
					j, _ := json.Marshal(zmsgs)
					w.WriteHeader(http.StatusCreated)
					w.Header().Set("Content-Type", "application/json")
					io.WriteString(w, string(j))
					return
				}
			}
		}
	} else {
		http.Error(w, "", 405)
		return
	}
}
