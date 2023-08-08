package handler

import (
	"GamePoolApi/backgroundmanager/enum/errorcode"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/reqireheader"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/zaplog"
	"net/http"
)

/////////////////////////////
////    轉調CQ9類型的共用結構/function
/////////////////////////////

// 文件標準HttpResponse結構
type BaseHttpResponse struct {
	Data   interface{}          `json:"data"`   //資料給予的地方
	Status MaHttpResponseStatus `json:"status"` //狀態欄
}

// 序列化BaseHttpResponse
func (res *BaseHttpResponse) ToString() string {
	data := serializer.JsonMarshal(res.Status.TraceCode, res)
	return string(data)
}

// 自訂Request Header
type BaseSelfDefine struct {
	TraceCode   string `json:"tracecode"`   //追蹤碼
	RequestTime string `json:"requesttime"` //http request time
}

// 文件標準HttpResponse.Status
type MaHttpResponseStatus struct {
	Code       string `json:"code"`                 //狀態碼
	Message    string `json:"message"`              //狀態訊息
	DateTime   string `json:"dateTime"`             //回應時間
	TraceCode  string `json:"tracecode"`            //追蹤碼
	Latency    string `json:"latency,omitempty"`    //latency
	WalletType string `json:"wallettype,omitempty"` //錢包類別transfer=轉帳錢包，single=單一錢包，ce=虛擬幣錢包
}

// 封裝Http Repsonse
func getErrorHttpResponse(code string, requestTime, traceCode, specErrorMessage string, data interface{}) BaseHttpResponse {
	resp := BaseHttpResponse{
		Data: data,
		Status: MaHttpResponseStatus{
			Code:      code,
			Message:   specErrorMessage,
			DateTime:  requestTime,
			TraceCode: traceCode,
		},
	}

	//log response body
	zaplog.Infow(logResponse, innertrace.FunctionNode, thirdparty.GetErrorHttpResponse, innertrace.TraceNode, traceCode, innertrace.DataNode, resp.ToString())

	return resp
}

// 取出request header的Authorization
func getAuthorizationFromRequest(r *http.Request) (authorization string) {
	return r.Header.Get(reqireheader.Authorization)
}

// 取出自訂request header的traceCode
func getTraceCodeFromRequest(r *http.Request) (traceCode string) {
	return r.Header.Get(innertrace.TraceNode)
}

// 取出自訂request header的requestTime
func getRequestTimeFromRequest(r *http.Request) (requestTime string) {
	return r.Header.Get(innertrace.RequestTimeNode)
}

// 封裝輸出錯誤的格式
func responseError(code errorcode.ErrorCode, requestTime, traceCode, specErrorMessage string) []byte {
	errorCode := string(code)
	resp := getErrorHttpResponse(errorCode, requestTime, traceCode, specErrorMessage, nil)
	byteResp := serializer.JsonMarshal(traceCode, resp)
	return byteResp
}

/////////////////////////////
////    純後台類型的共用結構/function
/////////////////////////////

// 後台API Http Response
type BackendBaseHttpResponse struct {
	Count  *int                      `json:"Count,omitempty"` //資料筆數
	Data   interface{}               `json:"Data"`            //資料給予的地方
	Status BackendHttpResponseStatus `json:"Status"`          //狀態
}

// 後台API HttpResponse.Status
type BackendHttpResponseStatus struct {
	Code      string `json:"Code"`      //狀態碼
	Message   string `json:"Message"`   //訊息
	Timestamp int64  `json:"Timestamp"` //時間戳
	//TOCHECK:目前後台API response沒有帶TraceCode,後續應考慮加上
	//TraceCode string `json:"tracecode"` //追蹤碼
}

// 序列化BackendBaseHttpResponse
func (resp *BackendBaseHttpResponse) ToString(traceCode string) string {
	data := serializer.JsonMarshal(traceCode, resp)
	return string(data)
}

// 封裝Backend Error Http Repsonse
func getBackendErrorHttpResponse(code string, traceCode, specErrorMessage string, data interface{}) BackendBaseHttpResponse {
	resp := BackendBaseHttpResponse{
		Data: data,
		Status: BackendHttpResponseStatus{
			Code:      code,
			Message:   specErrorMessage,
			Timestamp: mtime.UtcNow().Unix(),
			//TOCHECK:目前後台API response沒有帶TraceCode,後續應考慮加上
			//TraceCode: traceCode,
		},
	}

	//log response body
	zaplog.Infow(logResponse, innertrace.FunctionNode, thirdparty.GetBackendErrorHttpResponse, innertrace.TraceNode, traceCode, innertrace.DataNode, resp.ToString(traceCode))

	return resp
}

// 封裝輸出Backend Api錯誤的格式
func backendError(code errorcode.ErrorCode, traceCode, specErrorMessage string) []byte {
	errorCode := string(code)
	resp := getBackendErrorHttpResponse(errorCode, traceCode, specErrorMessage, nil)
	byteResp := serializer.JsonMarshal(traceCode, resp)
	return byteResp
}

/////////////////////////////
////    共用function
/////////////////////////////

// response回寫
func writeHttpResponse(w http.ResponseWriter, traceCode, errorCode string, data []byte) {
	//add response header
	w.Header().Add(reqireheader.ContentType, reqireheader.Json)
	w.Header().Add(reqireheader.TransferEncoding, reqireheader.Chunked)
	//需要這一行刷新header
	w.WriteHeader(http.StatusOK)

	//write response
	w.Write(data)
}
