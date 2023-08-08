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

//----------Probability Game Flow----------

/////////////////////////////
////    遊戲回合下注
/////////////////////////////

// 遊戲回合下注Request
type GameBetRequest struct {
	BaseSelfDefine        //自訂headers
	UserId         string `json:"id" validate:"gt=0"`                                         //玩家id
	GameToken      string `json:"gametoken" validate:"gt=0"`                                  //玩家遊戲token
	GameCode       string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	Round          string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Amount         string `json:"amount" validate:"amt"`                                      //下注金額※最大長度為12位數，及小數點後4位
	MtCode         string `json:"mtcode" validate:"mtcode"`                                   //交易代碼
	BetTime        string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //下注時間(UTC-4)
}

// 序列化GameBetRequest
func (req *GameBetRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 遊戲回合下注 Response
type GameBetResponse struct {
	Data   GameBetResponseData  `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// 遊戲回合下注 Response data
type GameBetResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary	Probability Game Bet (遊戲回合下注)
//	@Tags		Probability Game Flow
//	@Accept		x-www-form-urlencoded
//	@Param		id			formData	string	true	"玩家id"
//	@Param		gametoken	formData	string	true	"玩家遊戲token"
//	@Param		gamecode	formData	string	true	"遊戲代號"
//	@Param		round		formData	string	true	"遊戲回合編號"
//	@Param		amount		formData	float64	true	"下注金額※最大長度為12位數，及小數點後4位"
//	@Param		mtcode		formData	string	true	"交易代碼"
//	@Param		datetime	formData	string	true	"下注時間(UTC-4)"
//	@Success	200			{object}	GameBetResponse
//	@Router		/gamepool/CC/game/bet [post]
//
//	@Security	Bearer
func GameBet(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := GameBetRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.GameBet, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.UserId = r.FormValue("id")
	request.GameToken = r.FormValue("gametoken")
	request.GameCode = r.FormValue("gamecode")
	request.Round = r.FormValue("round")
	request.Amount = r.FormValue("amount")
	request.MtCode = r.FormValue("mtcode")
	request.BetTime = r.FormValue("datetime")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.GameBet, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.ProbabilityGameBet(traceCode, request.UserId, request.GameToken, request.GameCode, request.Round, request.Amount, request.MtCode, request.BetTime)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded ProbabilityGameBet failure!")
		return errResp
	}

	//轉換data
	data := GameBetResponse{
		Data: GameBetResponseData(cq9Data.Data),
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
