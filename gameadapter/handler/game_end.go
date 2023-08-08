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
////    遊戲回合結算
/////////////////////////////

// 遊戲回合結算Request
type GameEndRequest struct {
	BaseSelfDefine             //自訂headers
	UserId              string `json:"id" validate:"gt=0"`                                         //玩家id
	GameToken           string `json:"gametoken" validate:"gt=0"`                                  //玩家遊戲token
	GameCode            string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	Round               string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Jackpot             string `json:"jackpot" validate:"omitempty,amt"`                           //jackpot金額
	JackpotType         string `json:"jackpottype"`                                                //jackpot種類
	JackpotContriubtion string `json:"jackpotcontriubtion" validate:"omitempty,farray"`            //jackot貢獻值
	Detail              string `json:"detail"`                                                     //傳送此Slot game 的FreeGame, LuckyDraw, Bonus 的場次，使用string 傳送進來 "{"FreeGame": 10,"LuckyDraw": 2,"Bonus": 5}"
	EndTime             string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //結算時間(UTC-4)
	Rtp                 string `json:"rtp" validate:"omitempty,farray"`                            //殺率
	SingleRowBet        string `json:"singlerowbet" validate:"omitempty,boolean"`                  //true or false: 是否為再旋轉形成的注單
	Ip                  string `json:"ip"`                                                         //客端ip
}

// 序列化GameEndRequest
func (req *GameEndRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 遊戲回合結算 Response
type GameEndResponse struct {
	Data   GameEndResponseData  `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// 遊戲回合結算 Response data
type GameEndResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary	Probability Game End (遊戲回合結算)
//	@Tags		Probability Game Flow
//	@Accept		x-www-form-urlencoded
//	@Param		id					formData	string		true	"玩家id"
//	@Param		gametoken			formData	string		true	"玩家遊戲token"
//	@Param		gamecode			formData	string		true	"遊戲代號"
//	@Param		round				formData	string		true	"遊戲回合編號"
//	@Param		jackpot				formData	float64		false	"jackpot金額"
//	@Param		jackpottype			formData	string		false	"jackpot種類"
//	@Param		jackpotcontriubtion	formData	[]float64	false	"jackot貢獻值"
//	@Param		detail				formData	string		false	"傳送此Slot	game的FreeGame,LuckyDraw,Bonus的場次"
//	@Param		datetime			formData	string		true	"結算時間(UTC-4)"
//	@Param		rtp					formData	[]float64	false	"殺率"
//	@Param		singlerowbet		formData	bool		false	"是否為再旋轉形成的注單"
//	@Param		ip					formData	string		false	"客端ip"
//	@Success	200					{object}	GameEndResponse
//	@Router		/gamepool/CC/game/end [post]
//
//	@Security	Bearer
func GameEnd(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := GameEndRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.GameEnd, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.UserId = r.FormValue("id")
	request.GameToken = r.FormValue("gametoken")
	request.GameCode = r.FormValue("gamecode")
	request.Round = r.FormValue("round")
	request.EndTime = r.FormValue("datetime")
	//非必須欄位
	request.Jackpot = r.FormValue("jackpot")
	request.JackpotType = r.FormValue("jackpottype")
	request.JackpotContriubtion = r.FormValue("jackpotcontriubtion")
	request.Detail = r.FormValue("detail")
	request.Rtp = r.FormValue("rtp")
	request.SingleRowBet = r.FormValue("usecard")
	request.Ip = r.FormValue("ip")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.GameEnd, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.ProbabilityGameEnd(traceCode, request.UserId, request.GameToken, request.GameCode, request.Round, request.Jackpot, request.JackpotType, request.JackpotContriubtion, request.Detail, request.EndTime, request.Rtp, request.SingleRowBet, request.Ip)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded ProbabilityGameEnd failure!")
		return errResp
	}

	//轉換data
	data := GameEndResponse{
		Data: GameEndResponseData(cq9Data.Data),
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
