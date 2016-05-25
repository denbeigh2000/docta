package checkers

import (
	"github.com/denbeigh2000/docta"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestData struct {
	Endpoint string
	Method   string
	Data     []byte
}

func (r RequestData) Request() (*http.Response, error) {
	var data *bytes.Buffer
	if strings.ToLower(r.Method) == "post" {
		data = bytes.NewBuffer(r.Data)
	}

	req, err := http.NewRequest(r.Method, r.Endpoint, data)
	if err != nil {
		return nil, fmt.Errorf("Failure making HTTP request: %s")
	}

	return http.DefaultClient.Do(req)
}

type HTTPChecker struct {
	YellowContains string
	RedContains    string

	RequestData RequestData
}

func (c HTTPChecker) Check() docta.HealthState {
	resp, err := c.RequestData.Request()
	if err != nil {
		return docta.HealthState{docta.Red, err.Error()}
	}

	if resp.StatusCode != 200 {
		return docta.HealthState{docta.Red, fmt.Sprintf("Non-200 status code: %v", resp.StatusCode)}
	}

	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return docta.HealthState{docta.Red, err.Error()}
	}

	strBody := string(reqBody)
	if strings.Contains(strBody, c.RedContains) {
		return docta.HealthState{docta.Red, fmt.Sprintf("Response body contains forbidden term: %v", c.RedContains)}
	}

	if strings.Contains(strBody, c.YellowContains) {
		return docta.HealthState{docta.Yellow, fmt.Sprintf("Response body contains forbidden term: %v", c.YellowContains)}
	}

	return docta.HealthState{docta.Green, docta.DefaultInfo}
}
