package request

import (
	"fmt"
	"strconv"

	"bitbucket.com/LippoDigitalOVO/cashbackworker/cberror"
	"bitbucket.com/LippoDigitalOVO/cashbackworker/logger"
)

type httpLogger struct {
	req *HttpRequest
	res *httpResponse
}

func newHTTPLogger(req *HttpRequest, res *httpResponse) (log *httpLogger) {
	log = &httpLogger{}
	log.req = req
	log.res = res

	return
}

func (httpLogger *httpLogger) logData() {
	log := logger.GetLogger()
	reqBodyStringified := dumpMap("", httpLogger.req.Body)
	finalLog := "Url: " + httpLogger.req.Url + ", " +
		"ReqBody: " + reqBodyStringified + ", " +
		"StatusCode: " + strconv.Itoa(httpLogger.res.StatusCode) + ",  " +
		"ResBody: " + httpLogger.res.ResponseBody + ",  " +
		"Err: "
	log.Infof(cberror.HttpLoggerInfo, finalLog, httpLogger.res.Err)
}

func dumpMap(space string, m map[string]interface{}) (result string) {
	final := "{ "
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			final += fmt.Sprintf("{ \"%v\": ", k)
			final += dumpMap(space+"  ", mv)
			final += fmt.Sprintf("} ")
		} else {
			final += fmt.Sprintf("%v %v : %v ", space, k, v)
		}
	}
	return final + " }"
}
