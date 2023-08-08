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
////    遊戲回合贏分
/////////////////////////////

// 遊戲回合贏分Request
type GameWinRequest struct {
	BaseSelfDefine        //自訂headers
	UserId         string `json:"id" validate:"gt=0"`                                         //玩家id
	GameToken      string `json:"gametoken" validate:"gt=0"`                                  //玩家遊戲token
	GameCode       string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	Round          string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Amount         string `json:"amount" validate:"amt"`                                      //贏分金額※最大長度為12位數，及小數點後4位
	MtCode         string `json:"mtcode" validate:"mtcode"`                                   //交易代碼
	WinTime        string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //贏分時間(UTC-4)
	CardWin        string `json:"cardwin" validate:"omitempty,amt"`                           //派彩加成※最大長度為12位數，及小數點後4位
	UseCard        string `json:"usecard" validate:"omitempty,boolean"`                       //是否派彩加成
}

// 序列化GameWinRequest
func (req *GameWinRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 遊戲回合贏分 Response
type GameWinResponse struct {
	Data   GameWinResponseData  `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// 遊戲回合贏分 Response data
type GameWinResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary		Probability Game Win (遊戲回合贏分)
//	@Tags			Probability Game Flow
//
//	@Description	win 金額需要包含jackpot金額
//
//	@Accept			x-www-form-urlencoded
//	@Param			id			formData	string	true	"玩家id"
//	@Param			gametoken	formData	string	true	"玩家遊戲token"
//	@Param			gamecode	formData	string	true	"遊戲代號"
//	@Param			round		formData	string	true	"遊戲回合編號"
//	@Param			amount		formData	float64	true	"贏分金額※最大長度為12位數，及小數點後4位"
//	@Param			mtcode		formData	string	true	"交易代碼"
//	@Param			datetime	formData	string	true	"贏分時間(UTC-4)"
//	@Param			cardwin		formData	float64	false	"派彩加成※最大長度為12位數，及小數點後4位"
//	@Param			usecard		formData	bool	false	"是否派彩加成"
//	@Success		200			{object}	GameWinResponse
//	@Router			/gamepool/CC/game/win [post]
//
//	@Security		Bearer
func GameWin(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := GameWinRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.GameWin, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.UserId = r.FormValue("id")
	request.GameToken = r.FormValue("gametoken")
	request.GameCode = r.FormValue("gamecode")
	request.Round = r.FormValue("round")
	request.Amount = r.FormValue("amount")
	request.MtCode = r.FormValue("mtcode")
	request.WinTime = r.FormValue("datetime")
	//非必須欄位
	request.CardWin = r.FormValue("cardwin")
	request.UseCard = r.FormValue("usecard")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.GameWin, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.ProbabilityGameWin(traceCode, request.UserId, request.GameToken, request.GameCode, request.Round, request.Amount, request.MtCode, request.WinTime, request.CardWin, request.UseCard)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded ProbabilityGameWin failure!")
		return errResp
	}

	//轉換data
	data := GameWinResponse{
		Data: GameWinResponseData(cq9Data.Data),
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
