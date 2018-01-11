package main

import (
    "net/http"
    "bytes"
    "text/template"
)

func hold(acc string, amnt int) error {
	client := &http.Client{}

    // request with xml soap data
    reqXml, err := holdRequest("34444")
    if err != nil {
        println(err.Error)
		return err
    }
	req, err := http.NewRequest("POST", bankConfig.finacleApi, bytes.NewBuffer([]byte(reqXml)))
	if err != nil {
        println(err.Error)
		return err
	}

    // headers
	req.Header.Add("SOAPAction", `"http://ws.cdyne.com/WeatherWS/GetCityWeatherByZIP"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")

    // send request
    resp, err := client.Do(req)
	if err != nil {
        println(err.Error)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
        println("invalid response")
        // TODO
        return nil
	}

    // parse response and take account hold status

    return nil
}

func release(acc string, amnt int) {

}

func transfer(from string, to string, amount int) {

}

func holdRequest(postalCode string)(string, error) {
    type RequestParam struct {
        PostalCode string
    }

    reqTem := `
    <?xml version="1.0" encoding="utf-8">
    <soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">;
        <soap:Body>
            <GetCityWeatherByZIP xmlns="http://ws.cdyne.com/WeatherWS/">
	            <ZIP>{{.PostalCode}}</ZIP>
	        </GetCityWeatherByZIP>
        </soap:Body>
    </soap:Envelope>
    `
    p := RequestParam{PostalCode: postalCode}
    t, err := template.New("HoldRequest").Parse(reqTem)
    if err != nil {
        println(err.Error())
        return "", err
    }

    var buf bytes.Buffer
    err = t.Execute(&buf, p)
    if err != nil {
        println(err.Error())
        return "", err
    }

    return buf.String(), nil
}

func releaseRequest() {

}

func trasferRequest() {

}
