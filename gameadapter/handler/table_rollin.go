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

//----------Table Flow----------

/////////////////////////////
////    遊戲錢包轉至個人錢包
/////////////////////////////

// 遊戲錢包轉至個人錢包Request
type TableRollinRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken" validate:"gt=0"`                                  //玩家遊戲token
	UserId         string `json:"id" validate:"gt=0"`                                         //玩家id
	MtCode         string `json:"mtcode" validate:"mtcode"`                                   //交易代碼
	Round          string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Amount         string `json:"amount" validate:"amt"`                                      //結算金額※最大長度為12位數，及小數點後4位
	Bet            string `json:"bet" validate:"amt"`                                         //bet 金額
	Win            string `json:"win" validate:"amt"`                                         //win 金額(可為負值) ※含jackpot、不含抽水
	Rake           string `json:"rake" validate:"amt"`                                        //抽水金額
	GameRole       string `json:"gamerole" validate:"oneof=banker player"`                    //玩家角色為庄家(banker) or 閒家(player)
	BankerType     string `json:"bankertype" validate:"oneof=pc human"`                       //對戰玩家是否有真人[pc|human] pc：對戰玩家沒有真人 human：對戰玩家有真人 ※此欄位為牌桌遊戲使用，非牌桌遊戲此欄位值為空字串 ※如果玩家不支持上庄，只存在與系统對玩。則bankertype 為 PC
	RollinTime     string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //結算時間(UTC-4)
	GameCode       string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	RoomFee        string `json:"roomfee" validate:"omitempty,amt"`                           //開房費用
	ValidBet       string `json:"validbet" validate:"amt"`                                    //有效投注
	RoundNumber    string `json:"roundnumber"`                                                //局號
	TableType      string `json:"tabletype" validate:"omitempty,oneof=1 4"`                   //真人遊戲類別 1：百家樂 4：龍虎
	TableId        string `json:"tableid"`                                                    //桌號
	BetType        string `json:"bettype"`                                                    //下注玩法 ※真人參數說明https://hackmd.io/tbT46brhReC4Aq5mz5HM1w?view#bettype
	GameResult     string `json:"gameresult"`                                                 //遊戲結果 百家樂範例:{“points”:[5,6],“cards”:[{“poker”:“S2”,“tag”:2},{“poker”:“S8”,“tag”:1},{“poker”:“S3”,“tag”:2},{“poker”:“C7”,“tag”:1},{“poker”:“C1”,“tag”:2}]} 龍虎範例:{“points”:[9,2],“cards”:[{“poker”:“H9”,“tag”:1},{“poker”:“H2”,“tag”:2}]} ※真人參數說明https://hackmd.io/tbT46brhReC4Aq5mz5HM1w?view#bettype
	Ip             string `json:"ip"`                                                         //客端ip
	CardWin        string `json:"cardwin" validate:"omitempty,amt"`                           //派彩加成※最大長度為12位數，及小數點後4位
	UseCard        string `json:"usecard" validate:"omitempty,boolean"`                       //是否派彩加成
	LiveBetDetail  string `json:"livebetdetail"`                                              //視訊下注明細 範例:[{“bettype”:“2”,“win”:440,“winlose”:220,“bet”:220,“validbet”:220,“odds”:1,“matchresult”:“2”}] ※真人參數說明https://hackmd.io/tbT46brhReC4Aq5mz5HM1w?view#bettype
}

// 序列化TableRollinRequest
func (req *TableRollinRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 遊戲錢包轉至個人錢包 Response
type TableRollinResponse struct {
	Data   TableRollinResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus    `json:"status"` //狀態欄
}

// 遊戲錢包轉至個人錢包 Response data
type TableRollinResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary		Table Rollin Version 2 (遊戲錢包轉至個人錢包)
//	@Tags			Table Flow
//	@Description	Rollin amount = rollout amount + win - rake - roomfee
//	@Accept			x-www-form-urlencoded
//	@Param			gametoken		formData	string	true	"玩家遊戲token"
//	@Param			id				formData	string	true	"玩家id"
//	@Param			mtcode			formData	string	true	"交易代碼"
//	@Param			round			formData	string	true	"遊戲回合編號"
//	@Param			amount			formData	float64	true	"結算金額※最大長度為12位數，及小數點後4位"
//	@Param			bet				formData	float64	true	"bet金額"
//	@Param			win				formData	float64	true	"win 金額(可為負值) ※含jackpot、不含抽水"
//	@Param			rake			formData	float64	true	"抽水金額"
//	@Param			gamerole		formData	string	true	"玩家角色為庄家(banker) or 閒家(player)"
//	@Param			bankertype		formData	string	true	"對戰玩家是否有真人[pc|human] pc：對戰玩家沒有真人 human：對戰玩家有真人 ※此欄位為牌桌遊戲使用，非牌桌遊戲此欄位值為空字串"
//	@Param			datetime		formData	string	true	"結算時間(UTC-4)"
//	@Param			gamecode		formData	string	true	"遊戲代號"
//	@Param			roomfee			formData	float64	false	"開房費用"
//	@Param			validbet		formData	float64	true	"有效投注"
//	@Param			roundnumber		formData	string	false	"局號"
//	@Param			tabletype		formData	string	false	"真人遊戲類別 1：百家樂 4：龍虎"
//	@Param			tableid			formData	string	false	"桌號"
//	@Param			bettype			formData	[]int	false	"下注玩法※真人參數說明https://hackmd.io/tbT46brhReC4Aq5mz5HM1w?view#bettype"
//	@Param			gameresult		formData	string	false	"遊戲結果※真人參數說明https://hackmd.io/tbT46brhReC4Aq5mz5HM1w?view#bettype"
//	@Param			ip				formData	string	false	"客端ip"
//	@Param			cardwin			formData	float64	false	"派彩加成※最大長度為12位數，及小數點後4位"
//	@Param			usecard			formData	bool	false	"是否派彩加成"
//	@Param			livebetdetail	formData	string	false	"視訊下注明細※真人參數說明https://hackmd.io/tbT46brhReC4Aq5mz5HM1w?view#bettype"
//	@Success		200				{object}	TableRollinResponse
//	@Router			/gamepool/v2/CC/table/rollin [post]
//
//	@Security		Bearer
func TableRollin(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := TableRollinRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.TableRollin, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.GameToken = r.FormValue("gametoken")
	request.UserId = r.FormValue("id")
	request.MtCode = r.FormValue("mtcode")
	request.Round = r.FormValue("round")
	request.Amount = r.FormValue("amount")
	request.Bet = r.FormValue("bet")
	request.Win = r.FormValue("win")
	request.Rake = r.FormValue("rake")
	request.GameRole = r.FormValue("gamerole")
	request.BankerType = r.FormValue("bankertype")
	request.RollinTime = r.FormValue("datetime")
	request.GameCode = r.FormValue("gamecode")
	request.ValidBet = r.FormValue("validbet")
	//非必須欄位
	request.RoomFee = r.FormValue("roomfee")
	request.RoundNumber = r.FormValue("roundnumber")
	request.TableType = r.FormValue("tabletype")
	request.TableId = r.FormValue("tableid")
	request.BetType = r.FormValue("bettype")
	request.GameResult = r.FormValue("gameresult")
	request.Ip = r.FormValue("ip")
	request.CardWin = r.FormValue("cardwin")
	request.UseCard = r.FormValue("usecard")
	request.LiveBetDetail = r.FormValue("livebetdetail")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.TableRollin, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.TableRollinVersion2(traceCode, request.GameToken, request.UserId, request.MtCode, request.Round, request.Amount, request.Bet, request.Win, request.Rake, request.GameRole, request.BankerType, request.RollinTime, request.GameCode, request.RoomFee, request.ValidBet, request.RoundNumber, request.TableType, request.TableId, request.BetType, request.GameResult, request.Ip, request.CardWin, request.UseCard, request.LiveBetDetail)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded TableRollinVersion2 failure!")
		return errResp
	}

	//轉換data
	data := TableRollinResponse{
		Data: TableRollinResponseData(cq9Data.Data),
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
