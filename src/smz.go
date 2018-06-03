package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Smz struct {
	Appcode    string `json:"appcode"`
	Password   string `json:"password"`
	Message    string `json:"message"`
	MobileList string `json:"mobilelist"`
}

func send(zode string, zaddr string) error {
	// marshel senz
	msg := "iGift confirmation code " + zode
	smz := Smz{
		Appcode:    smzConfig.appcode,
		Password:   smzConfig.password,
		Message:    msg,
		MobileList: zaddr,
	}
	j, _ := json.Marshal(smz)

	println("sending sms to: " + smzConfig.api + " data: " + string(j))

	req, err := http.NewRequest("POST", smzConfig.api, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	// send to sms api
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("Invalid response")
	}

	// read response
	b, _ := ioutil.ReadAll(resp.Body)
	println(string(b))

	return nil
}
