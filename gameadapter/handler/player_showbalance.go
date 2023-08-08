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

//----------Player----------

/////////////////////////////
////    查詢餘額
/////////////////////////////

// 查詢餘額Request
type ShowBalanceRequest struct {
	BaseSelfDefine        //自訂headers
	UserId         string `json:"id" validate:"gt=0"` //玩家id
	GameCode       string `json:"gamecode"`           //遊戲代號
}

// 序列化ShowBalanceRequest
func (req *ShowBalanceRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 查詢餘額 Response
type ShowBalanceResponse struct {
	Data   ShowBalanceResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus    `json:"status"` //狀態欄
}

// 查詢餘額 Response data
type ShowBalanceResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// TOCHECK:文件的玩家id欄位是url path parameter,但原代碼跟前端接法是query string,如果不經過轉調API直接用同樣方式大CQ9會失敗
//
//	@Summary	Balance (查詢餘額)
//	@Tags		Player
//	@Param		id			query		string	true	"玩家id"	default(5995123703a236000175b842)
//	@Param		gamecode	query		string	false	"遊戲代號"	default(CC01)
//	@Success	200			{object}	ShowBalanceResponse
//	@Router		/gamepool/CC/player/balance [get]
//
//	@Security	Bearer
func ShowBalance(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := ShowBalanceRequest{}

	//文件id是讀url parameter,但目前前端傳過來是query string
	// // read id from url
	// vars := mux.Vars(r)
	// request.UserId = vars["id"]

	//read request query string
	request.UserId = r.URL.Query().Get("id")
	request.GameCode = r.URL.Query().Get("gamecode")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.ShowBalance, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.Balance(traceCode, request.UserId, request.GameCode)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded Balance failure!")
		return errResp
	}

	//轉換data
	data := ShowBalanceResponse{
		Data: ShowBalanceResponseData(cq9Data.Data),
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
