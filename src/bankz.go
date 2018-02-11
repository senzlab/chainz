package main

import (
    "os"
    "net/http"
    "bytes"
    "text/template"
	"io/ioutil"
    "path/filepath"
    "errors"
)

type LienAdd struct {
    Account string
    Module string
    Amount string
    Currency string
    Reason string
}

type LienMod struct {
    Id string
    Account string
    Module string
    Amount string
    Currency string
}

func lienAdd(acc string, amnt string) (string,error) {
	client := &http.Client{}

    // request with xml soap data
    reqXml, err := lienAddReq(acc, amnt)
    if err != nil {
        println(err.Error)
		return "", err
    }
	req, err := http.NewRequest("POST", finacleConfig.api, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
        println(err.Error)
		return "", err
	}

    // headers
	req.Header.Add("SOAPAction", finacleConfig.lienAddAction)
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("Accept", "text/xml")

    // send request
    resp, err := client.Do(req)
	if err != nil {
        println(err.Error)
		return "", err
	}
	defer resp.Body.Close()

    println(resp.StatusCode)
	if resp.StatusCode != 200 {
        println("invalid response")
        return "", errors.New("Invalid response")
	}

    // parse response and take account hold status
    resXml, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        println(err.Error)
		return "", err
	}
    println(string(resXml))

    return "323232", nil
}

func lienMod(acc string, lienId string)error {
	client := &http.Client{}

    // request with xml soap data
    reqXml, err := lienModReq(acc, lienId)
    if err != nil {
        println(err.Error)
		return err
    }
	req, err := http.NewRequest("POST", finacleConfig.api, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
        println(err.Error)
		return err
	}

    // headers
	req.Header.Add("SOAPAction", finacleConfig.lienModAction)
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("Accept", "text/xml")

    // send request
    resp, err := client.Do(req)
	if err != nil {
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

func lienAddReq(account string, amount string)(string, error) {
    // format template path
    cwd, _ := os.Getwd()
    tp := filepath.Join(cwd, "./template/lienadd.xml")
    println(tp)

    // template from file
    t, err := template.ParseFiles(tp)
    if err != nil {
        println(err.Error())
        return "", err
    }

    // lienadd params
    la := LienAdd{}
    la.Account = account
    la.Amount = amount
    la.Currency = "LKR"

    // parse template
    var buf bytes.Buffer
    err = t.Execute(&buf, la)
    if err != nil {
        println(err.Error())
        return "", err
    }

    return buf.String(), nil
}

func lienModReq(account string, lienId string)(string, error) {
    // format template path
    cwd, _ := os.Getwd()
    tp := filepath.Join(cwd, "./template/lienmod.xml")
    println(tp)

    // template from file
    t, err := template.ParseFiles(tp)
    if err != nil {
        println(err.Error())
        return "", err
    }

    // lienadd params
    lm := LienMod{}
    lm.Account = account
    lm.Id = lienId
    lm.Amount = "0"
    lm.Currency = "LKR"

    // parse template
    var buf bytes.Buffer
    err = t.Execute(&buf, lm)
    if err != nil {
        println(err.Error())
        return "", err
    }

    return buf.String(), nil
}
