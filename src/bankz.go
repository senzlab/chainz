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

type FundTrans struct {
    FromAcc string
    ToAcc string
    Amount string
    Commission string
    Memo string
    Date string
}

func doFundTrans(fromAcc string, toAcc string, amount string)error {
	client := &http.Client{}

    // request with xml soap data
    reqXml, err := fundTransReq(fromAcc, toAcc, amount)
    if err != nil {
        println(err.Error)
		return err
    }
    println(reqXml)

	req, err := http.NewRequest("POST", transConfig.api, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
        println(err.Error)
		return err
	}

    // headers
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

func fundTransReq(fromAcc string, toAcc string, amount string)(string, error) {
    // format template path
    cwd, _ := os.Getwd()
    tp := filepath.Join(cwd, "./template/ftrans.xml")
    println(tp)

    // template from file
    t, err := template.ParseFiles(tp)
    if err != nil {
        println(err.Error())
        return "", err
    }

    // trans params
    ft := FundTrans{}
    ft.FromAcc = fromAcc
    ft.ToAcc = toAcc
    ft.Amount = amount
    ft.Commission = transConfig.commission
    ft.Memo = ""
    ft.Date = ""

    // parse template
    var buf bytes.Buffer
    err = t.Execute(&buf, ft)
    if err != nil {
        println(err.Error())
        return "", err
    }

    return buf.String(), nil
}
