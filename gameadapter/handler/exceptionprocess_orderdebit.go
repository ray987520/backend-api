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

//----------Game Result Excepition Processing----------

/////////////////////////////
////    注單補扣款
/////////////////////////////

// 注單補扣款Request
type OrderDebitRequest struct {
	BaseSelfDefine        //自訂headers
	UserId         string `json:"id" validate:"gt=0"`                                         //玩家id
	GameCode       string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	MtCode         string `json:"mtcode" validate:"mtcode"`                                   //交易代碼
	Round          string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Amount         string `json:"amount" validate:"amt"`                                      //補扣款金額※最大長度為12位數，及小數點後4位
	DebitTime      string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //補扣款時間
}

// 序列化OrderDebitRequest
func (req *OrderDebitRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 注單補扣款 Response
type OrderDebitResponse struct {
	Data   OrderDebitResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus   `json:"status"` //狀態欄
}

// 注單補扣款 Response data
type OrderDebitResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary	Order Debit (注單補扣款)
//	@Tags		Game Result Excepition Processing
//	@Accept		x-www-form-urlencoded
//	@Param		id			formData	string	true	"玩家id"
//	@Param		gamecode	formData	string	true	"遊戲代號"
//	@Param		mtcode		formData	string	true	"交易代碼"
//	@Param		round		formData	string	true	"遊戲回合編號"
//	@Param		amount		formData	float64	true	"補扣款金額※最大長度為12位數，及小數點後4位"
//	@Param		datetime	formData	string	true	"補扣款時間"
//	@Success	200			{object}	OrderDebitResponse
//	@Router		/gamepool/CC/order/debit [post]
//
//	@Security	Bearer
func OrderDebit(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := OrderDebitRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.OrderDebit, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.UserId = r.FormValue("id")
	request.GameCode = r.FormValue("gamecode")
	request.MtCode = r.FormValue("mtcode")
	request.Round = r.FormValue("round")
	request.Amount = r.FormValue("amount")
	request.DebitTime = r.FormValue("datetime")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.OrderDebit, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.OrderDebit(traceCode, request.UserId, request.GameCode, request.MtCode, request.Round, request.Amount, request.DebitTime)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded OrderDebit failure!")
		return errResp
	}

	//轉換data
	data := OrderDebitResponse{
		Data: OrderDebitResponseData(cq9Data.Data),
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
