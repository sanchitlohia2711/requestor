package request

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

type HttpRequest struct {
	Url             string
	Method          string
	Timeout         int
	Body            map[string]interface{}
	Raw             []byte
	BodyContentType string
	Headers         map[string]string
	Retries         int
	LogData         bool
}

func ExecuteHttpRequest(httpRequest *HttpRequest) (httpResponse *httpResponse, err error) {
	httpRequest.setDefaults()
	err = httpRequest.setBody()
	if err != nil {
		return
	}
	httpResponse, err = httpRequest.Execute()
	return
}

func (httpRequest *HttpRequest) validate() (err error) {
	//Check Url
	if httpRequest.Url == "" {
		err = fmt.Errorf("Url cannot be empty")
		return
	}

	//Check Method
	if httpRequest.Method == "" {
		err = fmt.Errorf("Request Method cannot be empty")
		return
	}
	return
}

func (httpRequest *HttpRequest) setDefaults() {
	if httpRequest.Retries == 0 {
		httpRequest.Retries = 1
	}
}

func (httpRequest *HttpRequest) setBody() (err error) {
	if (httpRequest.Method == "POST" || httpRequest.Method == "PUT") && httpRequest.BodyContentType == "form" {
		raw, e := json.Marshal(httpRequest.Body)
		if e != nil {
			err = fmt.Errorf("Cannot marshal body")
			return
		}
		httpRequest.Raw = raw
	}
	return
}

//Execute : execute request
func (httpRequest *HttpRequest) Execute() (httpResponse *httpResponse, err error) {
	goReq := gorequest.New()

	goReq = goReq.CustomMethod(httpRequest.Method, httpRequest.Url)
	goReq = goReq.Timeout(time.Duration(httpRequest.Timeout) * time.Second)
	goReq = goReq.Type(httpRequest.BodyContentType)
	for k, v := range httpRequest.Headers {
		goReq.Set(k, v)
	}

	httpResponse, err = httpRequest.requestWithRetry(goReq)
	return
}

func (httpRequest *HttpRequest) requestWithRetry(goReq *gorequest.SuperAgent) (httpResponse *httpResponse, err error) {
	var res gorequest.Response
	var body string
	var errs []error
	for i := 0; i < httpRequest.Retries; i++ {
		if httpRequest.Method == "POST" || httpRequest.Method == "PUT" {
			res, body, errs = goReq.Send(string(httpRequest.Raw)).End()
		}
		if httpRequest.Method == "GET" {
			res, body, errs = goReq.Get(httpRequest.Url).End()
		}
		httpResponse = newHTTPResponse(res, body, errs)
		logger := newHTTPLogger(httpRequest, httpResponse)
		logger.logData()
		if !httpResponse.IsFailure() {
			break
		}
	}
	return
}
