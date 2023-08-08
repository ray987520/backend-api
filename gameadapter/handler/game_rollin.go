package handler

import (
	"GamePoolApi/common/database"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/str"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	cid "GamePoolApi/gameadapter/enum/controller"
	"GamePoolApi/gameadapter/enum/errorcode"
	"fmt"
	"math"
	"net/http"
)

//----------Probability Game Flow----------

/////////////////////////////
////    遊戲錢包轉至個人錢包
/////////////////////////////

// 遊戲錢包轉至個人錢包Request
type RollinRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken" validate:"gt=0"`                                  //玩家遊戲token
	UserId         string `json:"id" validate:"gt=0"`                                         //玩家id
	MtCode         string `json:"mtcode" validate:"mtcode"`                                   //交易代碼
	Round          string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Amount         string `json:"amount" validate:"amt"`                                      //結算金額※最大長度為12位數，及小數點後4位
	Bet            string `json:"bet" validate:"amt"`                                         //bet 金額
	Win            string `json:"win" validate:"amt"`                                         //win 金額 (含jackpot)
	Jackpot        string `json:"jackpot" validate:"omitempty,amt"`                           //jackpot 金額
	Item           string `json:"item" validate:"omitempty,amt"`                              //購買武器 金額
	Reward         string `json:"reward" validate:"omitempty,amt"`                            //遊戲獎金 金額
	RollinTime     string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //結算時間(UTC-4)
	GameCode       string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	Ip             string `json:"ip"`                                                         //客端ip
	CardWin        string `json:"cardwin" validate:"omitempty,amt"`                           //派彩加成※最大長度為12位數，及小數點後4位
	UseCard        string `json:"usecard" validate:"omitempty,boolean"`                       //是否派彩加成
}

// 序列化RollinRequest
func (req *RollinRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 遊戲錢包轉至個人錢包 Response
type RollinResponse struct {
	Data   RollinResponseData   `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// 遊戲錢包轉至個人錢包 Response data
type RollinResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary		Probability Game Rollin (遊戲錢包轉至個人錢包)
//	@Tags			Probability Game Flow
//
//	@Description	*win 金額需要包含jackpot金額
//	@Description	*bet 金額需包含購買武器金額
//	@Description	*Rollin amount = rollout amount - bet + win
//
//	@Accept			x-www-form-urlencoded
//	@Param			gametoken	formData	string	true	"玩家遊戲token"
//	@Param			id			formData	string	true	"玩家id"
//	@Param			mtcode		formData	string	true	"交易代碼"
//	@Param			round		formData	string	true	"遊戲回合編號"
//	@Param			amount		formData	float64	true	"結算金額※最大長度為12位數，及小數點後4位"
//	@Param			bet			formData	float64	true	"bet金額"
//	@Param			win			formData	float64	true	"win 金額 (含jackpot)"
//	@Param			jackpot		formData	float64	false	"jackpot金額"
//	@Param			item		formData	float64	false	"購買武器金額"
//	@Param			reward		formData	float64	false	"遊戲獎金金額"
//	@Param			datetime	formData	string	true	"結算時間"
//	@Param			gamecode	formData	string	true	"遊戲代號"
//	@Param			ip			formData	string	false	"客端ip"
//	@Param			cardwin		formData	float64	false	"派彩加成※最大長度為12位數，及小數點後4位"
//	@Param			usecard		formData	bool	false	"是否派彩加成"
//	@Success		200			{object}	RollinResponse
//	@Router			/gamepool/CC/rollin [post]
//
//	@Security		Bearer
func Rollin(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := RollinRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.Rollin, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.GameToken = r.FormValue("gametoken")
	request.UserId = r.FormValue("id")
	request.MtCode = r.FormValue("mtcode")
	request.Round = r.FormValue("round")
	request.Amount = r.FormValue("amount")
	request.RollinTime = r.FormValue("datetime")
	request.GameCode = r.FormValue("gamecode")
	request.Bet = r.FormValue("bet")
	request.Win = r.FormValue("win")
	//非必須欄位
	request.Jackpot = r.FormValue("jackpot")
	request.Item = r.FormValue("item")
	request.Reward = r.FormValue("reward")
	request.Ip = r.FormValue("ip")
	request.CardWin = r.FormValue("cardwin")
	request.UseCard = r.FormValue("usecard")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.Rollin, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉換輸贏
	temp, isOK := str.ParseFloat64(traceCode, request.Win)
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, fmt.Sprintf("parse win float64 failure! win:%s", request.Win))
		return errResp
	}
	win := int64(math.Round(temp * 10000.0))
	temp, isOK = str.ParseFloat64(traceCode, request.Bet)
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, fmt.Sprintf("parse bet float64 failure! bet:%s", request.Bet))
		return errResp
	}
	bet := int64(math.Round(temp * 10000.0))
	winLose := win - bet

	//轉換會員編號
	memberId := database.PlatformMemberIDChangeMemberID(traceCode, request.UserId)
	//轉換會員編號失敗輸出錯誤
	if memberId == -1 {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "db PlatformMemberIDChangeMemberID failure!")
		return errResp
	}

	//轉換遊戲代碼
	gameId := database.GameCodeChangeGameID(traceCode, request.GameCode)
	//轉換遊戲代碼失敗輸出錯誤
	if gameId == -1 {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "db GameCodeChangeGameID failure!")
		return errResp
	}

	//轉換時間
	payoutTime, isOK := mtime.ParseTime(traceCode, mtime.ApiTimeFormat, request.RollinTime)
	//轉換時間失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "transfer datetime error!")
		return errResp
	}
	payoutTimeStr := mtime.TimeString(payoutTime, mtime.DbTimeFormat)

	//db add game result
	resultCode := database.AddGameResult(traceCode, memberId, gameId, request.Round, bet, winLose, win, payoutTimeStr)
	switch resultCode {
	case 0: //繼續執行,0=成功
		break
	default:
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, fmt.Sprintf("db AddGameResult error! ResultCode:%d", resultCode))
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.ProbabilityGameRollin(traceCode, request.GameToken, request.UserId, request.MtCode, request.Round, request.Amount, request.Bet, request.Win, request.Jackpot, request.Item, request.Reward, request.RollinTime, request.GameCode, request.Ip, request.CardWin, request.UseCard)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded ProbabilityGameRollin failure!")
		return errResp
	}

	//轉換data
	data := RollinResponse{
		Data: RollinResponseData(cq9Data.Data),
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
