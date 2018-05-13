package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type FundTrans struct {
	FromAcc       string
	FromBankCode  string
	ToAcc         string
	ToBankCode    string
	Amount        string
	Commission    string
	CommissionAcc string
	Memo          string
	Date          string
}

type AccInq struct {
	Account string
}

func doAccVerify(acc string, statusVerify bool) error {
	client := &http.Client{}

	// request with xml soap data
	reqXml, err := verifyAccReq(acc)
	if err != nil {
		println(err.Error)
		return err
	}
	println(transConfig.api)
	println(reqXml)

	// TODO remove this
	return nil

	req, err := http.NewRequest("POST", transConfig.api, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
		println("error create request")
		println(err.Error)
		return err
	}

	// headers
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("Accept", "text/xml")

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
	resStr := string(resXml)
	println(resStr)

	if strings.Contains(resStr, "<AcctStatus>A") {
		// account active
		if statusVerify {
			return nil
		} else {
			if strings.Contains(resStr, "<ModeOfOperation>SNG") || strings.Contains(resStr, "<ModeOfOperation>ANY") {
				// inq done
				return nil
			}
		}
	}

	return errors.New("Invalid account")
}

func doFundTrans(fromAcc string, toAcc string, amount string, commission string, memo string) error {
	client := &http.Client{}

	// request with xml soap data
	reqXml, err := fundTransReq(fromAcc, toAcc, amount, commission, memo)
	if err != nil {
		println(err.Error)
		return err
	}
	println(transConfig.api)
	println(reqXml)

	// TODO remove this
	return nil

	req, err := http.NewRequest("POST", transConfig.api, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
		println("error create request")
		println(err.Error)
		return err
	}

	// headers
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("Accept", "text/xml")

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
	if !strings.Contains(string(resXml), "<ActionCode>000") {
		// trans done
		println("invalid response")
		return errors.New("Invalid response")
	}

	return nil
}

func doCeftTrans(toAcc string, bankCode string, amount string, memo string) error {
	client := &http.Client{}

	// request with xml soap data
	reqXml, err := ceftTransReq(toAcc, bankCode, amount, memo)
	if err != nil {
		println(err.Error)
		return err
	}
	println(transConfig.api)
	println(reqXml)

	// TODO remove this
	return nil

	req, err := http.NewRequest("POST", transConfig.api, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
		println("error create request")
		println(err.Error)
		return err
	}

	// headers
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("Accept", "text/xml")

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
	if !strings.Contains(string(resXml), "<TxnStatus>000") {
		// trans done
		println("invalid response")
		return errors.New("Invalid response")
	}

	return nil
}

func fundTransReq(fromAcc string, toAcc string, amount string, commission string, memo string) (string, error) {
	// format template path
	cwd, _ := os.Getwd()
	tp := filepath.Join(cwd, "./template/trans.xml")
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
	ft.Commission = commission
	ft.Memo = memo
	ft.Date = time.Now().Format("02/01/2006")

	// parse template
	var buf bytes.Buffer
	err = t.Execute(&buf, ft)
	if err != nil {
		println(err.Error())
		return "", err
	}

	return buf.String(), nil
}

func ceftTransReq(acc string, bankCode string, amount string, memo string) (string, error) {
	// format template path
	cwd, _ := os.Getwd()
	tp := filepath.Join(cwd, "./template/ceft-trans.xml")
	println(tp)

	// template from file
	t, err := template.ParseFiles(tp)
	if err != nil {
		println(err.Error())
		return "", err
	}

	// trans params
	ft := FundTrans{}
	ft.FromAcc = transConfig.account
	ft.FromBankCode = transConfig.bankCode
	ft.ToAcc = acc
	ft.ToBankCode = bankCode
	ft.Amount = amount
	ft.Commission = transConfig.ceftCommission
	ft.CommissionAcc = transConfig.commissionAccount
	ft.Date = time.Now().Format("02/01/2006")
	ft.Memo = memo

	// parse template
	var buf bytes.Buffer
	err = t.Execute(&buf, ft)
	if err != nil {
		println(err.Error())
		return "", err
	}

	return buf.String(), nil
}

func verifyAccReq(acc string) (string, error) {
	// format template path
	cwd, _ := os.Getwd()
	tp := filepath.Join(cwd, "./template/accinq.xml")
	println(tp)

	// template from file
	t, err := template.ParseFiles(tp)
	if err != nil {
		println(err.Error())
		return "", err
	}

	// trans params
	inq := AccInq{}
	inq.Account = acc

	// parse template
	var buf bytes.Buffer
	err = t.Execute(&buf, inq)
	if err != nil {
		println(err.Error())
		return "", err
	}

	return buf.String(), nil
}
