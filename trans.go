package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	api    = "http://192.168.125.93:7800/sd/iib/IIBFinacleIntegration"
	action = "http://192.168.125.93:7800/sd/iib/IIBFinacleIntegration"
)

func main() {
	trans()
}

func trans() error {
	client := &http.Client{}

	// request with xml soap data
	reqXml := req()
	println(reqXml)
	println("----1")

	req, err := http.NewRequest("POST", api, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
		println("----2")
		println(err.Error)
		return err
	}

	// headers
	//req.Header.Add("SOAPAction", action)
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("Accept", "text/xml")

	// send request
	resp, err := client.Do(req)
	if err != nil {
		println("----3")
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

func req() string {
	xml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:iib="http://www.sampath.lk/SD/IIBFinacleIntegration/">
   <soapenv:Header/>
   <soapenv:Body>
      <iib:DoTransferRequest>
         <APPCode>SVR</APPCode>
         <Controller>CMN</Controller>
         <CDCICode>C</CDCICode>
         <FromAccountNo>100105875594</FromAccountNo>
         <ToAccountNo>100105999635</ToAccountNo>
         <DTxnAmount>100</DTxnAmount>
         <DCommAmount></DCommAmount>
         <TransMemo>Promize Eranga</TransMemo>
         <ValueDate>06/03/2018</ValueDate>
         <FromCurrCode>LKR</FromCurrCode>
      </iib:DoTransferRequest>
   </soapenv:Body>
</soapenv:Envelope>
`
	return xml
}
