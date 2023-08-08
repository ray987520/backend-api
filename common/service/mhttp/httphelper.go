package mhttp

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"moul.io/http2curl"
)

/////////////////////////////
////    封裝的http共用服務
/////////////////////////////

const (
	logResponse                 = "log response body"
	newHttpRequestError         = "create new http request error:%v"
	newHttpRequestFormdataError = "create new http request formdata error:%v"
)

// 對url編碼避免特殊字元
func UrlEncode(data string) string {
	return url.QueryEscape(data)
}

// 對url已編碼特殊字元解碼
func UrlDecode(data string) (string, error) {
	decodeData, err := url.QueryUnescape(data)
	//url decode error,return error
	if err != nil {
		return "", err
	}
	return decodeData, nil
}

// 將http request內容轉成curl字串,方便log與復現
func HttpRequest2Curl(req *http.Request) (string, error) {
	curl, err := http2curl.GetCurlCommand(req)
	//gen curl string error, return error
	if err != nil {
		return "", err
	}

	return curl.String(), nil
}

// 調用第三方API,返回response body
func CallThirdApi(traceCode, httpMethod, requrl string, header map[string]string, formData map[string]string) (responseBody []byte, isOK bool) {
	var (
		req *http.Request
		err error
	)

	zaplog.Infow(innertrace.LogThirdApiParameters, innertrace.FunctionNode, thirdparty.CallThirdApi, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("httpMethod", httpMethod, "requrl", requrl, "header", header, "formData", formData))

	//判斷formData是否nil產生GET/非GET httprequest
	if formData != nil {
		/* MARK 記錄,如果需要使用multipart/form使用這種
		//set http request formData(multipart/form,C#用,跟header不符但看起來要用這種沒有encode才不會跳time format error)
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		for k, v := range formData {
			pw, err := w.CreateFormField(k)
			if err != nil {
				err = fmt.Errorf(newHttpRequestFormdataError, err)
				zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.CallThirdApi, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "k", k, "v", v))
				return nil, false
			}
			pw.Write([]byte(v))
		}
		*/
		//set http request formData(x-www-form-urlencoded)
		postData := url.Values{}
		for k, v := range formData {
			postData.Add(k, v)
		}

		req, err = http.NewRequest(httpMethod, requrl, strings.NewReader(postData.Encode()))
	} else {
		req, err = http.NewRequest(httpMethod, requrl, nil)
	}

	//new http request error,返回失敗
	if err != nil {
		err = fmt.Errorf(newHttpRequestError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.CallThirdApi, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "req", req))
		return nil, false
	}

	//set http request header
	for k, v := range header {
		req.Header.Set(k, v)
	}

	//log call thirdparty api request
	curl, err := HttpRequest2Curl(req)
	zaplog.Infow(innertrace.LogThirdApiRequest, innertrace.FunctionNode, thirdparty.CallThirdApi, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("curl", curl, innertrace.ErrorInfoNode, err))

	//call third api
	response, err := http.DefaultClient.Do(req)
	//調用異常,記錄並返回失敗
	if err != nil {
		err = fmt.Errorf(newHttpRequestError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.CallThirdApi, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err))
		return nil, false
	}

	//read response body
	defer response.Body.Close()
	responseBody, err = ioutil.ReadAll(response.Body)

	//log call third api response
	zaplog.Infow(innertrace.LogThirdApiResponse, innertrace.FunctionNode, thirdparty.CallThirdApi, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("responseBody", string(responseBody), "response.StatusCode", response.StatusCode))

	//read response body error,返回失敗
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.CallThirdApi, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err))
		return nil, false
	}

	return responseBody, isResponseSuccessStatusCode(response.StatusCode)
}

// 調用第三方API,返回response body
func CallThirdApiByBody(traceCode, httpMethod, requrl string, header map[string]string, body []byte) (responseBody []byte, isOK bool) {
	var (
		req *http.Request
		err error
	)

	zaplog.Infow(innertrace.LogThirdApiParameters, innertrace.FunctionNode, thirdparty.CallThirdApiByBody, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("httpMethod", httpMethod, "requrl", requrl, "header", header, "body", string(body)))

	req, err = http.NewRequest(httpMethod, requrl, bytes.NewReader(body))

	//new http request error,返回失敗
	if err != nil {
		err = fmt.Errorf(newHttpRequestError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.CallThirdApiByBody, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "req", req))
		return nil, false
	}

	//set http request header
	for k, v := range header {
		req.Header.Set(k, v)
	}

	//log call thirdparty api request
	curl, err := HttpRequest2Curl(req)
	zaplog.Infow(innertrace.LogThirdApiRequest, innertrace.FunctionNode, thirdparty.CallThirdApiByBody, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("curl", curl, innertrace.ErrorInfoNode, err))

	//call third api
	response, err := http.DefaultClient.Do(req)
	//調用異常,記錄並返回失敗
	if err != nil {
		err = fmt.Errorf(newHttpRequestError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.CallThirdApiByBody, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err))
		return nil, false
	}

	//read response body
	defer response.Body.Close()
	responseBody, err = ioutil.ReadAll(response.Body)

	//log call third api response
	zaplog.Infow(innertrace.LogThirdApiResponse, innertrace.FunctionNode, thirdparty.CallThirdApiByBody, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("responseBody", string(responseBody), "response.StatusCode", response.StatusCode))

	//read response body error,返回失敗
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.CallThirdApiByBody, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err))
		return nil, false
	}

	return responseBody, isResponseSuccessStatusCode(response.StatusCode)
}

// 判斷response statuscode是否是成功,200~299
func isResponseSuccessStatusCode(code int) bool {
	return code >= 200 && code <= 299
}

// 讀取http request body
func ReadHttpRequestBody(traceCode string, request *http.Request) ([]byte, bool) {
	defer request.Body.Close()
	requestBody, err := ioutil.ReadAll(request.Body)
	//讀取body失敗返回錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.ReadHttpRequestBody, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "requestBody", string(requestBody)))
		return nil, false
	}

	return requestBody, true
}

// 把struct轉成QueryString
func ToQueryString(traceCode string, data interface{}) (string, bool) {
	v, err := query.Values(data)
	//導入結構失敗返回錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.ToQueryString, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "data", data))
		return "", false
	}

	return v.Encode(), true
}
