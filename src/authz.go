package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type AuthParam struct {
	Username string
	Password string
}

func doAuth(user *User) error {
	client := &http.Client{}

	// request with xml soap data
	reqXml, err := authReq(user)
	if err != nil {
		println(err.Error)
		return err
	}
	println(authConfig.api)
	println(reqXml)

	// TODO remove this
	//return nil

	req, err := http.NewRequest("POST", authConfig.api, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
		println("error create request")
		println(err.Error)
		return err
	}

	// headers
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("Accept", "text/xml")
	req.Header.Add("SOAPAction", authConfig.action)

	// send request
	resp, err := client.Do(req)
	if err != nil {
		println("error call request")
		println(err.Error)
		return err
	}
	defer resp.Body.Close()

	println(resp.StatusCode)
	if resp.StatusCode != 200 {
		println("invalid response")
		return errors.New("Invalid response")
	}

	// parse response and take account hold status
	resXml, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error)
		return err
	}
	println(string(resXml))

	return nil
}

func authReq(user *User) (string, error) {
	// format template path
	cwd, _ := os.Getwd()
	tp := filepath.Join(cwd, "./template/auth.xml")
	println(tp)

	// template from file
	t, err := template.ParseFiles(tp)
	if err != nil {
		println(err.Error())
		return "", err
	}

	// auth param
	authParam := AuthParam{}
	authParam.Username = "eranga"
	authParam.Password = "1234"

	// parse template
	var buf bytes.Buffer
	err = t.Execute(&buf, authParam)
	if err != nil {
		println(err.Error())
		return "", err
	}

	return buf.String(), nil
}
