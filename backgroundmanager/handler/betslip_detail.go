package handler

import (
	cid "GamePoolApi/backgroundmanager/enum/controller"
	"GamePoolApi/backgroundmanager/enum/errorcode"
	"GamePoolApi/backgroundmanager/enum/gameresult"
	"GamePoolApi/backgroundmanager/enum/gametype"
	"GamePoolApi/backgroundmanager/enum/respmsg"
	"GamePoolApi/common/database"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/crypt"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/str"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	"encoding/json"
	"net/http"
	"strings"
)

//---------注單----------

/////////////////////////////
////    注單詳情
/////////////////////////////

// 注單詳情Request
type BetSlipDetailsRequest struct {
	BaseSelfDefine
	RoundId string `json:"RoundID" validate:"gt=0"` //局號
}

// 序列化BetSlipDetailsRequest
func (req *BetSlipDetailsRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 注單詳情Response
type BetSlipDetailsResponse struct {
	Data   BetSlipDetailsResponseData `json:"Data"`   //資料給予的地方
	Status BackendHttpResponseStatus  `json:"Status"` //狀態
}

// 注單詳情Response data
type BetSlipDetailsResponseData struct {
	StatusID         int                  `json:"StatusID"`         //狀態碼	0(正常)、1(會員不存在)、2(交易失敗)、3(寫賽果失敗)、4(寫log失敗)
	OwnerID          string               `json:"OwnerID"`          // 總代理編號
	ParentID         string               `json:"ParentID"`         // 代理編號
	Paccount         string               `json:"Paccount"`         // 代理帳號
	PlatformMemberID string               `json:"PlatformMemberID"` // 對方會員編號
	MemberAccount    string               `json:"MemberAccount"`    // 對方會員帳號
	GameCode         string               `json:"GameCode"`         // 遊戲代碼
	GameName         string               `json:"GameName"`         // 遊戲名稱
	Currency         string               `json:"Currency"`         // 幣別
	GameResult       GameResultDetailList `json:"GameResult"`       // 遊戲結果
	RoundID          string               `json:"RoundID"`          // 局號
	GameLog          json.RawMessage      `json:"GameLog"`          // 遊戲log
	EndTime          string               `json:"EndTime"`          // 成單時間
}

// 賽果詳情列表
type GameResultDetailList struct {
	GameResultDetails []GameResultDetail `json:"GameResultDetails"` // 賽果詳情
}

// 賽果詳情
type GameResultDetail struct {
	Time   string  `json:"Time"`   // 時間
	Action string  `json:"Action"` // 動作
	Amount float64 `json:"Amount"` // 金額
}

//	@Summary	注單詳情
//	@Tags		注單
//	@Accept		json
//	@Param		RoundID	query		string	true	"局號"	default(CCVY1h5pfvunb000694)
//	@Success	200		{object}	BetSlipDetailsResponse
//	@Router		/BetSlip/BetSlipDetails [get]
//
//	@Security	Bearer
func BetSlipDetails(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := BetSlipDetailsRequest{}

	//read query string
	request.RoundId = r.URL.Query().Get("RoundID")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.BetSlipDetails, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := backendError(errorcode.BackendError, traceCode, "bad request data!")
		return errResp
	}

	// 從jwt token取得使用者訊息
	claim := crypt.JwtValidAccessToken(traceCode, getAuthorizationFromRequest(r))
	userTimeZone := "-04:00"
	if claim != nil {
		userTimeZone = claim.TimeZone
	}
	timeZone, isOK := str.Atoi(traceCode, strings.Split(userTimeZone, ":")[0]) //jwt token時區是+08:00之類的字串,須轉為int使用
	//轉換userTimeZone失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "transfer userTimeZone error!")
		return errResp
	}

	gameLog, isOK := database.GameLogGet(traceCode, request.RoundId)
	//取出遊戲詳情失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "DB Platfrom_GameLogGet failure!")
		return errResp
	}

	//轉換遊戲類型編號 1:slot 2:fish
	gameType := ""
	switch gameLog.GameTypeID {
	case 1:
		gameType = gametype.Slot
	case 2:
		gameType = gametype.Fish
	default:
		errResp := backendError(errorcode.BackendError, traceCode, "DB Platfrom_GameLogGet error!")
		return errResp
	}

	//取CQ9遊戲細單客端網址
	_, cq9DetailTokenData, isOK := cq9.CreateDetailToken(traceCode, request.RoundId, gameLog.GameCode, gameType)
	//取CQ9遊戲細單客端網址失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "send cq9 detailtoken failure!")
		return errResp
	}

	//取回CQ9細單資訊
	_, cq9DetailData, isOK := cq9.DetailToken(traceCode, cq9DetailTokenData.Data.DetailToken)
	//取回CQ9細單資訊失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "send cq9 DetailTokenResult failure!")
		return errResp
	}

	//設置response
	//轉換格式
	betTime := ""
	if gameLog.BetTime != "" && gameLog.BetTime != mtime.DateTimeOffsetMinValue {
		dbBetTime, isOK := mtime.ParseToTimeZone(traceCode, mtime.ApiTimeFormat, gameLog.BetTime, timeZone)
		//轉換dbBetTime失敗輸出錯誤
		if !isOK {
			errResp := backendError(errorcode.BackendError, traceCode, "transfer dbBetTime error!")
			return errResp
		}

		betTime = mtime.TimeStringAndFillZero(dbBetTime, mtime.ApiTimeFormat)
	}
	payoutTime := ""
	if gameLog.PayoutTime != "" && gameLog.PayoutTime != mtime.DateTimeOffsetMinValue {
		dbPayoutTime, isOK := mtime.ParseToTimeZone(traceCode, mtime.ApiTimeFormat, gameLog.PayoutTime, timeZone)
		//轉換dbPayoutTime失敗輸出錯誤
		if !isOK {
			errResp := backendError(errorcode.BackendError, traceCode, "transfer dbPayoutTime error!")
			return errResp
		}

		payoutTime = mtime.TimeStringAndFillZero(dbPayoutTime, mtime.ApiTimeFormat)
	}
	endTime := ""
	if gameLog.EndTime != "" && gameLog.EndTime != mtime.DateTimeOffsetMinValue {
		dbEndTime, isOK := mtime.ParseToTimeZone(traceCode, mtime.ApiTimeFormat, gameLog.EndTime, timeZone)
		//轉換dbEndTime失敗輸出錯誤
		if !isOK {
			errResp := backendError(errorcode.BackendError, traceCode, "transfer dbEndTime error!")
			return errResp
		}

		endTime = mtime.TimeStringAndFillZero(dbEndTime, mtime.ApiTimeFormat)
	}
	gameDetail := []byte(nil)
	if gameLog.GameLog != "" {
		gameDetail = []byte(gameLog.GameLog)
	}

	//組成賽果詳情
	gameResultDetails := []GameResultDetail{}
	gameResultDetails = append(gameResultDetails, GameResultDetail{
		Time:   betTime,
		Action: gameresult.Bet,
		Amount: float64(gameLog.Bet) / 10000.0,
	})
	gameResultDetails = append(gameResultDetails, GameResultDetail{
		Time:   payoutTime,
		Action: gameresult.Win,
		Amount: float64(gameLog.Payout) / 10000.0,
	})

	data := BetSlipDetailsResponseData{
		OwnerID:          gameLog.OwnerID,
		ParentID:         gameLog.ParentID,
		Paccount:         cq9DetailData.Data.PAccount,
		PlatformMemberID: gameLog.PlatformMemberID,
		MemberAccount:    gameLog.MemberAccount,
		GameCode:         gameLog.GameCode,
		GameName:         gameLog.GameName,
		RoundID:          gameLog.RoundID,
		Currency:         gameLog.Currency,
		GameLog:          gameDetail,
		StatusID:         gameLog.StatusID,
		GameResult:       GameResultDetailList{GameResultDetails: gameResultDetails},
		EndTime:          endTime,
	}
	response := BetSlipDetailsResponse{
		Data: data,
		Status: BackendHttpResponseStatus{
			Code:      string(errorcode.Success),
			Message:   respmsg.Success,
			Timestamp: mtime.UtcNow().Unix(),
			//TraceCode: traceCode,
		},
	}
	byteResponse := serializer.JsonMarshal(traceCode, response)

	return byteResponse
}
