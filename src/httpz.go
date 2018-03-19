package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

	// router
	r := mux.NewRouter()
	r.HandleFunc("/promizes", promizes).Methods("POST")
	r.HandleFunc("/uzers", uzers).Methods("POST")

	// start server
	err := http.ListenAndServe(":7070", r)
	if err != nil {
		println(err.Error)
		os.Exit(1)
	}
}

func promizes(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var senzMsg SenzMsg
	json.Unmarshal(b, &senzMsg)

	// handle senz
	senz := parse(senzMsg.Msg)
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
		err := doFundTrans(trans.FromAccount, trans.ToAccount, trans.PromizeAmount, transConfig.commission)
		if err != nil {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}

		// create cheque
		// create trans
		createPromize(promize)
		createTrans(trans)

		// TODO handle create failures

		// msg to #from
		// message for #to
		fMsg := SenzMsg{
			Uid: senz.Attr["uid"],
			Msg: statusSenz("SUCCESS", senz.Attr["uid"], senz.Sender),
		}
		tMsg := SenzMsg{
			Uid: senz.Attr["uid"],
			Msg: promizeSenz(promize, senz.Sender, senz.Attr["to"], uid()),
		}

		// msgs
		var zmsgs []SenzMsg
		zmsgs = append(zmsgs, fMsg)
		zmsgs = append(zmsgs, tMsg)

		// success response
		successResponse(w, zmsgs)
		return
	} else {
		// this mean already transfered cheque
		// check for double spend
		if isDoubleSpend(senz.Sender, id) {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}

		// get cheque first
		promize, err := getPromize(senz.Attr["bnk"], id)
		if err != nil {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}

		// new trans
		trans := senzToTrans(&senz, promize)
		trans.FromAccount = transConfig.account
		trans.ToBank = senz.Attr["bnk"]
		trans.ToAccount = senz.Attr["acc"]
		trans.Type = "REDEEM"

		// call finacle to fund transfer
		err = doFundTrans(trans.FromAccount, trans.ToAccount, trans.PromizeAmount, transConfig.commission)
		if err != nil {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}

		// create trans
		createTrans(trans)

		// status to sender
		zmsg := SenzMsg{
			Uid: senz.Attr["uid"],
			Msg: statusSenz("SUCCESS", senz.Attr["uid"], senz.Sender),
		}

		// msgs
		var zmsgs []SenzMsg
		zmsgs = append(zmsgs, zmsg)

		// success response
		successResponse(w, zmsgs)
		return
	}
}

func uzers(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var senzMsg SenzMsg
	json.Unmarshal(b, &senzMsg)

	// handle senz
	senz := parse(senzMsg.Msg)
	if _, ok := senz.Attr["pubkey"]; ok {
		// create user
		user := senzToUser(&senz)
		err := createUser(user)
		if err != nil {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}

		// success response
		zmsg := SenzMsg{
			Uid: senz.Attr["uid"],
			Msg: statusSenz("SUCCESS", senz.Attr["uid"], senz.Sender),
		}
		var zmsgs []SenzMsg
		zmsgs = append(zmsgs, zmsg)

		// success response back
		successResponse(w, zmsgs)
		return
	}

	if acc, ok := senz.Attr["acc"]; ok {
		// add account
		// generate salt amount
		salt := randomSalt()

		// fund transfer salt amount
		err := doFundTrans(acc, transConfig.account, salt, "")
		if err != nil {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}

		// save salt
		// save account
		setUserSalt(salt, senz.Sender)
		setUserAccount(acc, senz.Sender)

		// return success response
		zmsg := SenzMsg{
			Uid: senz.Attr["uid"],
			Msg: statusSenz("SUCCESS", senz.Attr["uid"], senz.Sender),
		}
		var zmsgs []SenzMsg
		zmsgs = append(zmsgs, zmsg)

		// success response back
		successResponse(w, zmsgs)
		return
	}

	if salt, ok := senz.Attr["salt"]; ok {
		// confirm salt
		// get user and compare salt
		user, err := getUser(senz.Sender)
		if err != nil {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}
		if user.Salt != salt {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}

		// fund transfer salt amount
		err = doFundTrans(transConfig.account, user.Account, salt, "")
		if err != nil {
			errorResponse(w, senz.Attr["uid"], senz.Sender)
			return
		}

		// set state
		setUserState("VERIFIED", senz.Sender)

		// return success response
		zmsg := SenzMsg{
			Uid: senz.Attr["uid"],
			Msg: statusSenz("SUCCESS", senz.Attr["uid"], senz.Sender),
		}
		var zmsgs []SenzMsg
		zmsgs = append(zmsgs, zmsg)

		// success response back
		successResponse(w, zmsgs)
		return
	}
}

func errorResponse(w http.ResponseWriter, uid string, to string) {
	// marshel and return error
	senzMsg := SenzMsg{
		Uid: uid,
		Msg: statusSenz("ERROR", uid, to),
	}
	j, _ := json.Marshal(senzMsg)
	http.Error(w, string(j), 400)
}

func successResponse(w http.ResponseWriter, zmsgs []SenzMsg) {
	j, _ := json.Marshal(zmsgs)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(j))
}
