package request

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/parnurzeal/gorequest"
)

type httpResponse struct {
	ResponseBody string
	StatusCode   int
	Err          error
	Errs         []error
}

func newHTTPResponse(res gorequest.Response, body string, errs []error) (response *httpResponse) {
	response = &httpResponse{}
	if res != nil {
		response.ResponseBody = body
		response.StatusCode = res.StatusCode
	}
	if len(errs) > 0 {
		response.Err = errs[0]
		response.Errs = errs
	}
	return
}

func (res *httpResponse) IsFailure() (isFailure bool) {
	if res.Err != nil {
		isFailure = true
	}
	r, _ := regexp.Compile(`2\d\d`)
	if r.MatchString(strconv.Itoa(res.StatusCode)) == true {
		isFailure = false
		return
	}
	isFailure = true
	return
}

func (res *httpResponse) Is4xx() (is4xx bool) {
	r, _ := regexp.Compile(`4\d\d`)
	if r.MatchString(strconv.Itoa(res.StatusCode)) == true {
		is4xx = true
		return
	}
	is4xx = false
	return
}

func (res *httpResponse) Is5xx() (is4xx bool) {
	r, _ := regexp.Compile(`5\d\d`)
	if r.MatchString(strconv.Itoa(res.StatusCode)) == true {
		is4xx = true
		return
	}
	is4xx = false
	return
}

func (res *httpResponse) Is4Hundred() (is4xx bool) {
	r, _ := regexp.Compile(`4\d\d`)
	if r.MatchString(strconv.Itoa(res.StatusCode)) == true {
		is4xx = true
		return
	}
	is4xx = false
	return
}

func (res *httpResponse) IsTimeout() (isTimeout bool) {
	isTimeout = false
	for _, err := range res.Errs {
		if os.IsTimeout(err) {
			isTimeout = true
		}
	}
	return
}

func (res *httpResponse) ErrorResponse() (errorResponse error) {
	errorResponseString := "Status Code: " + strconv.Itoa(res.StatusCode) + "\n" +
		"Res Body: " + res.ResponseBody + "\n" +
		"Err: %v"
	errorResponse = fmt.Errorf(errorResponseString, res.Err)
	return
}
