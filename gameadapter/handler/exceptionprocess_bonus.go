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
////    派發紅利
/////////////////////////////

// 派發紅利Request
type BonusRequest struct {
	GameCode  string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	GameType  string `json:"gametype" validate:"gt=0"`                                   //遊戲類型
	Account   string `json:"account" validate:"acct"`                                    //帳號
	OwnerId   string `json:"ownerid" validate:"gt=0"`                                    //上層代理id
	ParentId  string `json:"parentid" validate:"gt=0"`                                   //代理ID
	UserId    string `json:"id" validate:"gt=0"`                                         //玩家id※此值為唯一值，請勿使用 account 替代 id
	Round     string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Amount    string `json:"amount" validate:"amt"`                                      //紅利金額
	BonusTime string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //派發紅利時間
}

// 派發紅利 Response
type BonusResponse struct {
	Data   BonusResponseData    `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// 派發紅利 Response data
type BonusResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary	Bonus (派發紅利)
//	@Tags		Game Result Excepition Processing
//	@Accept		x-www-form-urlencoded
//	@Param		gamecode	formData	string	true	"遊戲代號"
//	@Param		gametype	formData	string	true	"遊戲類型"
//	@Param		account		formData	string	true	"帳號"
//	@Param		ownerid		formData	string	true	"上層代理id"
//	@Param		parentid	formData	string	true	"代理ID"
//	@Param		id			formData	string	true	"玩家id"
//	@Param		round		formData	string	true	"遊戲回合編號"
//	@Param		amount		formData	float64	true	"紅利金額"
//	@Param		datetime	formData	string	true	"派發紅利時間"
//	@Success	200			{object}	BonusResponse
//	@Router		/gamepool/CC/order/bonus [post]
//
//	@Security	Bearer
func Bonus(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := BonusRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.Bonus, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.GameCode = r.FormValue("gamecode")
	request.GameType = r.FormValue("gametype")
	request.Account = r.FormValue("account")
	request.OwnerId = r.FormValue("ownerid")
	request.ParentId = r.FormValue("parentid")
	request.UserId = r.FormValue("id")
	request.Round = r.FormValue("round")
	request.Amount = r.FormValue("amount")
	request.BonusTime = r.FormValue("datetime")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.Bonus, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.Bonus(traceCode, request.GameCode, request.GameType, request.Account, request.OwnerId, request.ParentId, request.UserId, request.Round, request.Amount, request.BonusTime)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded Bonus failure!")
		return errResp
	}

	//轉換data
	data := BonusResponse{
		Data: BonusResponseData(cq9Data.Data),
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
