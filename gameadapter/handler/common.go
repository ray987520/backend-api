package handler

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/reqireheader"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/zaplog"
	"GamePoolApi/gameadapter/enum/errorcode"
	"net/http"
)

/////////////////////////////
////    共用http request/response結構
/////////////////////////////

// 文件標準HttpResponse結構
type BaseHttpResponse struct {
	Data   interface{}          `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
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
type GaHttpResponseStatus struct {
	Code       string `json:"code"`                 //狀態碼
	Message    string `json:"message"`              //狀態訊息
	DateTime   string `json:"dateTime"`             //回應時間
	TraceCode  string `json:"tracecode"`            //追蹤碼
	Latency    string `json:"latency,omitempty"`    //latency,TOCHECK:文件沒有但實際CQ9 response有
	WalletType string `json:"wallettype,omitempty"` //錢包類別transfer=轉帳錢包，single=單一錢包，ce=虛擬幣錢包
}

/////////////////////////////
////    共用封裝http method
/////////////////////////////

// 封裝錯誤的Http Repsonse
func getErrorHttpResponse(code string, requestTime, traceCode, specErrorMessage string, data interface{}) BaseHttpResponse {
	resp := BaseHttpResponse{
		Data: data,
		Status: GaHttpResponseStatus{
			Code:       code,
			Message:    specErrorMessage,
			DateTime:   requestTime,
			TraceCode:  traceCode,
			Latency:    "",
			WalletType: "",
		},
	}

	//log response body
	zaplog.Infow(logResponse, innertrace.FunctionNode, thirdparty.GetErrorHttpResponse, innertrace.TraceNode, traceCode, innertrace.DataNode, resp.ToString())

	return resp
}

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
