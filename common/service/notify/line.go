package notify

import (
	"GamePoolApi/common/enum/httpmethod"
	"GamePoolApi/common/enum/reqireheader"
	"net/http"
	"net/url"
	"strings"
)

/////////////////////////////
////    封裝的line警示服務
/////////////////////////////

// line警示實體結構
type LineNotify struct{}

const (
	notify_URL = "https://notify-api.line.me/api/notify" //line警示服務url
	sendError  = "send notify error"                     //傳送line警示失敗
)

var (
	accessToken       = ""        //line警示服務token,格式為"Bearer XXXXXXXXXXXXXXX"
	lineNotifyService *LineNotify // line警示實體
)

// 初始化line警示實體,傳入accessToken
func InitLineNotify(token string) {
	accessToken = token
	lineNotifyService = &LineNotify{}
}

// 實作io.Writer,可以加到zaplog同步傳送錯誤
func (ln *LineNotify) Write(message []byte) (n int, err error) {
	return Send(string(message))
}

// 傳送錯誤到line
func Send(message string) (n int, err error) {
	//設定http request
	header := map[string]string{}
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode //application/x-www-form-urlencoded
	header[reqireheader.Authorization] = accessToken
	formData := map[string]string{}
	formData["message"] = message

	var req *http.Request
	//set http request formData(x-www-form-urlencoded)
	postData := url.Values{}
	for k, v := range formData {
		postData.Add(k, v)
	}

	req, err = http.NewRequest(httpmethod.Post, notify_URL, strings.NewReader(postData.Encode()))

	//new http request error返回失敗
	if err != nil {
		return 400, err
	}

	//set http request header
	req.Header.Set(reqireheader.ContentType, reqireheader.FormUrlEncode) //application/x-www-form-urlencoded
	req.Header.Set(reqireheader.Authorization, accessToken)              //Bearer XXXXXXXXXXXXXXX

	//call third api
	response, err := http.DefaultClient.Do(req)
	//調用異常返回失敗
	if err != nil {
		return 400, err
	}

	return response.StatusCode, nil
}
