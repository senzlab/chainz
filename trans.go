package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	api    = "http://192.168.125.93:7800/sd/iib/IIBFinacleIntegration"
	action = "http://192.168.125.93:7800/sd/iib/IIBFinacleIntegration"
)

func main() {
	//trans()
	//date()
	//random()
	random()
	//y := random()
	//z := random()
	//a := fmt.Sprintf("%d.%d%d", x, y, z)
	//println(a)
}

func random() string {
	rand.Seed(time.Now().Unix())
	x := rand.Intn(2)
	y := rand.Intn(9-1) + 1
	z := rand.Intn(9)
	s := fmt.Sprintf("%d.%d%d", x, y, z)
	fmt.Println(s)

	return s
}

func date() {
	d := time.Now().Format("02/01/2006")
	println(d)
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

	if !strings.Contains(string(resXml), "<ActionCode>000") {
		// trans done
		println("invalid response----")
		return errors.New("Invalid response")
	} else {
		println("done res----")
	}

	return nil
}

func req() string {
	xml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:iib="http://www.sampath.lk/SD/IIBFinacleIntegration/">
   <soapenv:Header/>
   <soapenv:Body>
      <iib:DoTransferRequest>
         <APPCode>GFT</APPCode>
         <Controller>CMN</Controller>
         <CDCICode>C</CDCICode>
         <FromAccountNo>100105875594</FromAccountNo>
         <ToAccountNo>100105999635</ToAccountNo>
         <DTxnAmount>100</DTxnAmount>
         <DCommAmount></DCommAmount>
         <TransMemo>Promize Eranga</TransMemo>
         <ValueDate>19/03/2018</ValueDate>
         <FromCurrCode>LKR</FromCurrCode>
      </iib:DoTransferRequest>
   </soapenv:Body>
</soapenv:Envelope>
`
	return xml
}

func sreq() string {
	xml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:iib="http://www.sampath.lk/SD/IIBFinacleIntegration/">
   <soapenv:Header/>
   <soapenv:Body>
      <iib:DoInterBankTransferRequest>
         <!--Optional:-->
         <APPCode>VSW</APPCode>
         <!--Optional:-->
         <Controller>CMN</Controller>
         <!--Optional:-->
         <CDCICode>F</CDCICode>
         <!--Optional:-->
         <TerminalID></TerminalID>
         <!--Optional:-->
         <CardNo></CardNo>
         <FromAccNo>900100000801</FromAccNo>
         <!--Optional:-->
         <FromAccType></FromAccType>
         <!--Optional:-->
         <FromAccBankCode>7278</FromAccBankCode>
         <!--Optional:-->
         <FromAccBranchCode></FromAccBranchCode>
         <ToAccNo>100178099707</ToAccNo>
         <!--Optional:-->
         <ToAccName>Lakshan</ToAccName>
         <!--Optional:-->
         <ToAccType></ToAccType>
         <ToAccBankCode>7719</ToAccBankCode>
         <!--Optional:-->
         <ToAccBranchCode>017</ToAccBranchCode>
         <TxnAmount>100</TxnAmount>
         <CommAmount>20</CommAmount>
         <!--Optional:-->
         <CommAccount>900108020041</CommAccount>
         <!--Optional:-->
         <TranRemarks></TranRemarks>
         <!--Optional:-->
         <ValueDate>05/04/2018</ValueDate>
         <!--Optional:-->
         <SlipsCode></SlipsCode>
         <!--Optional:-->
         <DrCurrencyCode></DrCurrencyCode>
         <!--Optional:-->
         <ChannelType></ChannelType>
         <CEFTFlag></CEFTFlag>
         <FromAccName></FromAccName>
         <!--Optional:-->
         <Reference></Reference>
      </iib:DoInterBankTransferRequest>
   </soapenv:Body>
</soapenv:Envelope>
`
	return xml
}
