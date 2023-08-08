package handler

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	cid "GamePoolApi/gameadapter/enum/controller"
	"GamePoolApi/gameadapter/enum/errorcode"
	"net/http"
)

//---------接入輔助工具----------

/////////////////////////////
////    存款(將錢存至測試帳號)
/////////////////////////////

// 存款(將錢存至測試帳號)Request
type MoneyInRequest struct {
	BaseSelfDefine        //自訂headers
	Account        string `json:"account" validate:"acct"`                                                                         //玩家帳號
	Amount         string `json:"amount" validate:"amt"`                                                                           //金額
	Currency       string `json:"currency" validate:"oneof=CNY MYR THB RUB JPY KRW IDR USD IDR(K) VND(K) EUR SGD HKD INR MMK VND"` //使用幣別※目前支援CNY, MYR, THB, RUB, JPY, KRW, IDR, USD, IDR(K), VND(K), EUR, SGD, HKD, INR, MMK, VND
}

// 序列化MoneyInRequest
func (req *MoneyInRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 存款(將錢存至測試帳號) Response
type MoneyInResponse struct {
	Data   MoneyInResponseData  `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// 存款(將錢存至測試帳號) Response data
type MoneyInResponseData struct {
	Before    float64 `json:"before"`    //交易前餘額
	Balance   float64 `json:"balance"`   //交易後餘額
	Currency  string  `json:"currency"`  //幣別
	Tracecode float64 `json:"tracecode"` //追蹤碼
}

//	@Summary		存款(將錢存至測試帳號)
//	@Tags			接入輔助工具
//	@Description	account 與 amount 需要是正確的才能存錢至測試帳號
//	@Description	{相對應站點的URL}/peace/money/in
//	@Accept			x-www-form-urlencoded
//	@Param			account		formData	string	true	"玩家帳號"																							default(test001)
//	@Param			amount		formData	float64	true	"金額"																							default(10.1)
//	@Param			currency	formData	string	true	"使用幣別※目前支援CNY, MYR, THB, RUB, JPY, KRW, IDR, USD, IDR(K), VND(K), EUR, SGD, HKD, INR, MMK, VND"	default(CNY)
//	@Success		200			{object}	MoneyInResponse
//	@Router			/peace/money/in [post]
//
//	@Security		Bearer
func MoneyIn(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := MoneyInRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.MoneyIn, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.Account = r.FormValue("account")
	request.Amount = r.FormValue("amount")
	request.Currency = r.FormValue("currency")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.MoneyIn, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.MoneyIn(traceCode, request.Account, request.Amount, request.Currency)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded MoneyIn failure!")
		return errResp
	}

	//轉換data
	data := MoneyInResponse{
		Data: MoneyInResponseData(cq9Data.Data),
		Status: GaHttpResponseStatus{
			Code:       cq9Data.Status.Code,
			Message:    cq9Data.Status.Message,
			DateTime:   requestTime,
			TraceCode:  cq9Data.Status.TraceCode,
			Latency:    cq9Data.Status.Latency,
			WalletType: cq9Data.Status.WalletType,
		},
	}
	byteResponse := serializer.JsonMarshal(traceCode, data)

	return byteResponse
}
