package cq9

import (
	"GamePoolApi/common/entity"
	"GamePoolApi/common/enum/currency"
	"GamePoolApi/common/enum/errorcode"
	"GamePoolApi/common/enum/httpmethod"
	"GamePoolApi/common/enum/reqireheader"
	iface "GamePoolApi/common/interface"
	"GamePoolApi/common/service/mhttp"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/serializer"
	"fmt"
)

/////////////////////////////
////    CQ9服務
/////////////////////////////

var (
	cq9GameHost  string //cq9 host
	cq9AuthToken string //cq9 token
	cq9GameHall  string //cq9廳號
	cq9TeamID    string //cq9團隊代號
)

// 初始化CQ9 service
func InitCQ9Service(cfg iface.IAppConfig) {
	cq9GameHost = cfg.GetCq9GameHost()
	cq9AuthToken = cfg.GetCq9AuthToken()
	cq9GameHall = cfg.GetCq9GameHall()
	cq9TeamID = cfg.GetCq9TeamID()
}

//---------GamePool API----------

// Auth (驗證 Game Token，並且回傳此玩家的所有設定資訊)
func Auth(traceCode string, gameToken string) ([]byte, entity.AuthPlayerResponse, bool) {
	var (
		data entity.AuthPlayerResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/player/auth", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/player/auth
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["gametoken"] = gameToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Balance (查詢餘額)
func Balance(traceCode string, userId, gameCode string) ([]byte, entity.ShowBalanceResponse, bool) {
	var (
		data entity.ShowBalanceResponse
		isOK bool
	)

	//設置request內容
	//gamecode不為空時傳送url encode
	if gameCode != "" {
		gameCode = mhttp.UrlEncode(gameCode)
	}
	reqUrl := fmt.Sprintf("%s/gamepool/%s/player/balance/%s?gamecode=%s", cq9GameHost, cq9GameHall, userId, gameCode) //{Config.Url}/gamepool/{Config.GameHall}/player/balance/{id}?gamecode={gamecode}
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, reqUrl, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Logout (登出)
func Logout(traceCode string, gameToken string) ([]byte, entity.CQ9BaseHttpResponse, bool) {
	var (
		data entity.CQ9BaseHttpResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/player/logout", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/player/logout
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["gametoken"] = gameToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Probability Game Bet (遊戲回合下注)
func ProbabilityGameBet(traceCode string, userId, gameToken, gameCode, round, amount, mtCode, betTime string) ([]byte, entity.GameBetResponse, bool) {
	var (
		data entity.GameBetResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/bet", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/bet
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["id"] = userId
	formData["gametoken"] = gameToken
	formData["gamecode"] = gameCode
	formData["round"] = round
	formData["amount"] = amount
	formData["mtcode"] = mtCode
	formData["datetime"] = betTime

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Probability Game Win (遊戲回合贏分)
func ProbabilityGameWin(traceCode string, userId, gameToken, gameCode, round, amount, mtCode, winTime, cardWin, useCard string) ([]byte, entity.GameWinResponse, bool) {
	var (
		data entity.GameWinResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/win", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/win
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["id"] = userId
	formData["gametoken"] = gameToken
	formData["gamecode"] = gameCode
	formData["round"] = round
	formData["amount"] = amount
	formData["mtcode"] = mtCode
	formData["datetime"] = winTime
	if cardWin != "" {
		formData["cardwin"] = cardWin
	}
	if useCard != "" {
		formData["usecard"] = useCard
	}

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Probability Game End (遊戲回合結算)
func ProbabilityGameEnd(traceCode string, userId, gameToken, gameCode, round, jackpot, jackpotType, jackpotContriubtion, detail, endTime, rtp, singleRowBet, ip string) ([]byte, entity.GameEndResponse, bool) {
	var (
		data entity.GameEndResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/end", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/end
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode

	formData := map[string]string{}
	formData["id"] = userId
	formData["gametoken"] = gameToken
	formData["gamecode"] = gameCode
	formData["round"] = round
	if jackpot != "" {
		formData["jackpot"] = jackpot
	}
	if jackpotType != "" {
		formData["jackpottype"] = jackpotType
	}
	if jackpotContriubtion != "" {
		formData["jackpotcontriubtion"] = jackpotContriubtion
	}
	if detail != "" {
		formData["detail"] = detail
	}
	formData["datetime"] = endTime
	if rtp != "" {
		formData["rtp"] = rtp
	}
	if singleRowBet != "" {
		formData["singlerowbet"] = singleRowBet
	}
	if ip != "" {
		formData["ip"] = ip
	}

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Probability Game Win Cancel (贏分取消)
func ProbabilityGameWinCancel(traceCode string, mtCode string) ([]byte, entity.GameWinCancelResponse, bool) {
	var (
		data entity.GameWinCancelResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/win/cancel", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/win/cancel
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode

	formData := map[string]string{}
	formData["mtcode"] = mtCode

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Probability Game Rollout (個人錢包轉至遊戲錢包)
func ProbabilityGameRollout(traceCode string, gameToken, userId, mtCode, round, amount, rolloutTime, gameCode, takeAll string) ([]byte, entity.RolloutResponse, bool) {
	var (
		data entity.RolloutResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/rollout", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/rollout
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["gametoken"] = gameToken
	formData["id"] = userId
	formData["mtcode"] = mtCode
	formData["round"] = round
	formData["amount"] = amount
	formData["datetime"] = rolloutTime
	formData["gamecode"] = gameCode
	if takeAll != "" {
		formData["takeall"] = takeAll
	}

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Probability Game Rollin (遊戲錢包轉至個人錢包)
func ProbabilityGameRollin(traceCode string, gameToken, userId, mtCode, round, amount, bet, win, jackpot, item, reward, rollinTime, gameCode, ip, cardWin, useCard string) ([]byte, entity.RollinResponse, bool) {
	var (
		data entity.RollinResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/rollin", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/rollin
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["gametoken"] = gameToken
	formData["id"] = userId
	formData["mtcode"] = mtCode
	formData["round"] = round
	formData["amount"] = amount
	formData["bet"] = bet
	formData["win"] = win
	if jackpot != "" {
		formData["jackpot"] = jackpot
	}
	if item != "" {
		formData["item"] = item
	}
	if reward != "" {
		formData["reward"] = reward
	}
	formData["datetime"] = rollinTime
	formData["gamecode"] = gameCode
	if ip != "" {
		formData["ip"] = ip
	}
	if cardWin != "" {
		formData["cardwin"] = cardWin
	}
	if useCard != "" {
		formData["usecard"] = useCard
	}

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Table Rollout (將金額從個人錢包轉至遊戲錢包)
func TableRollout(traceCode string, gameToken, userId, mtCode, round, amount, rolloutTime, gameCode, gameRole, bankerType, takeAll string) ([]byte, entity.TableRolloutResponse, bool) {
	var (
		data entity.TableRolloutResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/table/rollout", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/table/rollout
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["gametoken"] = gameToken
	formData["id"] = userId
	formData["mtcode"] = mtCode
	formData["round"] = round
	formData["amount"] = amount
	formData["datetime"] = rolloutTime
	formData["gamecode"] = gameCode
	formData["gamerole"] = gameRole
	formData["bankertype"] = bankerType
	if takeAll != "" {
		formData["takeall"] = takeAll
	}

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Table Rollin Version 2 (遊戲錢包轉至個人錢包)
func TableRollinVersion2(traceCode string, gameToken, userId, mtCode, round, amount, bet, win, rake, gameRole, bankerType, rollinTime, gameCode, roomFee, validBet, roundNumber, tableType, tableId, betType, gameResult, ip, cardWin, useCard, liveBetDetail string) ([]byte, entity.TableRollinResponse, bool) {
	var (
		data entity.TableRollinResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/v2/%s/table/rollin", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/v2/{Config.GameHall}/table/rollin
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["gametoken"] = gameToken
	formData["id"] = userId
	formData["mtcode"] = mtCode
	formData["round"] = round
	formData["amount"] = amount
	formData["bet"] = bet
	formData["win"] = win
	formData["rake"] = rake
	formData["gamerole"] = gameRole
	formData["bankertype"] = bankerType
	formData["datetime"] = rollinTime
	formData["gamecode"] = gameCode
	if roomFee != "" {
		formData["roomfee"] = roomFee
	}
	formData["validbet"] = validBet
	if roundNumber != "" {
		formData["roundnumber"] = roundNumber
	}
	if tableType != "" {
		formData["tabletype"] = tableType
	}
	if tableId != "" {
		formData["tableid"] = tableId
	}
	if betType != "" {
		formData["bettype"] = betType
	}
	if gameResult != "" {
		formData["gameresult"] = gameResult
	}
	if ip != "" {
		formData["ip"] = ip
	}
	if cardWin != "" {
		formData["cardwin"] = cardWin
	}
	if useCard != "" {
		formData["usecard"] = useCard
	}
	if liveBetDetail != "" {
		formData["livebetdetail"] = liveBetDetail
	}

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Probability / Table Game Refund
func TableGameRefund(traceCode string, mtCode, userId string) ([]byte, entity.TableGameRefundResponse, bool) {
	var (
		data entity.TableGameRefundResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/bet/refund", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/bet/refund
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["mtcode"] = mtCode
	if userId != "" {
		formData["userid"] = userId
	}

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Detail Token (驗證細單Token正確性並取回相關細單資訊)
func DetailToken(traceCode string, token string) ([]byte, entity.DetailTokenResponse, bool) {
	var (
		data entity.DetailTokenResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/detailtoken", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/detailtoken
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["token"] = token

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	//如果返回status.Code不是0(成功)返回null
	if data.Status.Code != string(errorcode.Success) {
		return resp, data, false
	}

	return resp, data, true
}

// Round Check (查詢未完成注單 (RoundID 必須符合格式才能查詢到))
func RoundCheck(traceCode string, fromDate, toDate string) ([]byte, entity.RoundCheckResponse, bool) {
	var (
		data entity.RoundCheckResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/roundcheck?fromdate=%s&todate=%s", cq9GameHost, cq9GameHall, fromDate, toDate) //{Config.Url}/gamepool/{Config.GameHall}/game/roundcheck
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, url, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Round Detail (查詢注單內容)
func RoundDetail(traceCode string, indexId string) ([]byte, entity.RoundDetailResponse, bool) {
	var (
		data entity.RoundDetailResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/rounddetail?indexid=%s", cq9GameHost, cq9GameHall, indexId) //{Config.Url}/gamepool/{Config.GameHall}/game/rounddetail
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, url, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Get Order (查詢注單內容)
func GetOrder(traceCode string, indexId string) ([]byte, entity.GetOrderResponse, bool) {
	var (
		data entity.GetOrderResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/getorder?indexid=%s", cq9GameHost, cq9GameHall, indexId) //{Config.Url}/gamepool/{Config.GameHall}/game/getorder
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, url, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Recover Gametoken by time range (當重送end時，遇到 error code 9 Invalid GameToken,即可使用這個 API 重新激活相對應時間區間內的 gametoken)
func RecoverGametokenByTimeRange(traceCode string, fromDate, toDate string) ([]byte, entity.RecoverGametokenByTimeRangeResponse, bool) {
	var (
		data entity.RecoverGametokenByTimeRangeResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/generateroundcache", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/generateroundcach
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["fromdate"] = fromDate
	formData["todate"] = toDate

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Recover Gametoken by RoundID (當重送 End 時，遇到 error code 9 Invalid GameToken.，即可使用這個API重新激活相對應 RoundID 的 gametoken)
func RecoverGametokenByRoundID(traceCode string, indexId string) ([]byte, entity.RecoverGametokenByRoundIDResponse, bool) {
	var (
		data entity.RecoverGametokenByRoundIDResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/generateoneroundcache", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/generateoneroundcache
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["indexid"] = indexId

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Order Debit (注單補扣款)
func OrderDebit(traceCode string, userId, gameCode, mtCode, round, amount, debitTime string) ([]byte, entity.OrderDebitResponse, bool) {
	var (
		data entity.OrderDebitResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/order/debit", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/order/debit
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["id"] = userId
	formData["gamecode"] = gameCode
	formData["mtcode"] = mtCode
	formData["round"] = round
	formData["amount"] = amount
	formData["datetime"] = debitTime

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Order Credit (注單補款)
func OrderCredit(traceCode string, userId, gameCode, mtCode, round, amount, creditTime string) ([]byte, entity.OrderCreditResponse, bool) {
	var (
		data entity.OrderCreditResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/order/credit", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/order/credit
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["id"] = userId
	formData["gamecode"] = gameCode
	formData["mtcode"] = mtCode
	formData["round"] = round
	formData["amount"] = amount
	formData["datetime"] = creditTime

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Bonus (派發紅利)
func Bonus(traceCode string, gameCode, gameType, account, ownerId, parentId, userId, round, amount, bonusTime string) ([]byte, entity.BonusResponse, bool) {
	var (
		data entity.BonusResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/order/bonus", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/order/bonus
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["gamecode"] = gameCode
	formData["gametype"] = gameType
	formData["account"] = account
	formData["ownerid"] = ownerId
	formData["parentid"] = parentId
	formData["id"] = userId
	formData["round"] = round
	formData["amount"] = amount
	formData["datetime"] = bonusTime

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Currency List (幣別列表)
func CurrencyList(traceCode string) ([]byte, entity.CurrencyListResponse, bool) {
	var (
		data entity.CurrencyListResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/currency", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/currency
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, url, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Player Order URL(取得玩家注單網址)
func PlayerOrder(traceCode string, gameToken string) ([]byte, entity.PlayerOrderResponse, bool) {
	var (
		data entity.PlayerOrderResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/playerorder?gametoken=%s", cq9GameHost, cq9GameHall, gameToken) //{Config.Url}/gamepool/{Config.GameHall}/game/playerorder?gametoken=
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, url, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// Promotion(檢查目前是否有活動列表)
func Promotion(traceCode string, gameToken string) ([]byte, entity.PromotionResponse, bool) {
	var (
		data entity.PromotionResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/promo?gametoken=%s", cq9GameHost, cq9GameHall, gameToken) //{Config.Url}/gamepool/{Config.GameHall}/game/promo?gametoken=
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, url, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9第二層解析,按文件Promotion回傳的可能是array或nil,但只要status.code=0都是正常,所以status.code!=0返回錯誤
	if data.Status.Code != "0" {
		return resp, data, false
	}

	return resp, data, true
}

// Promotion Link(產生活動連結)
func PromotionLink(traceCode string, gameToken, promotionId string) ([]byte, entity.PromotionLinkResponse, bool) {
	var (
		data entity.PromotionLinkResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/gamepool/%s/game/promo/link", cq9GameHost, cq9GameHall) //{Config.Url}/gamepool/{Config.GameHall}/game/promo/link
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["gametoken"] = gameToken
	formData["promoid"] = promotionId

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

//---------輔助工具----------

// 創建gametoken(用於測試遊戲客端網址)
func CreateGameToken(traceCode string, account, gameType, gameCode, currency string) ([]byte, entity.CreateGameTokenResponse, bool) {
	var (
		data entity.CreateGameTokenResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/peace/gametoken?account=%s&gametype=%s&gamecode=%s", cq9GameHost, account, gameType, gameCode) //{Config.Url}/peace/gametoken?
	if currency != "" {
		url += fmt.Sprintf("&currency=%s", currency)
	}
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, url, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// 創建detailtoken(遊戲細單客端網址)
func CreateDetailToken(traceCode string, roundId, gameCode, gameType string) ([]byte, entity.CreateDetailTokenResponse, bool) {
	var (
		data entity.CreateDetailTokenResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/peace/detailtoken?roundid=%s&gamecode=%s&gametype=%s", cq9GameHost, roundId, gameCode, gameType) //{Config.Url}/peace/detailtoken?
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, url, header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	//如果返回status.Code不是0(成功)返回null
	if data.Status.Code != string(errorcode.Success) {
		return resp, data, false
	}

	return resp, data, true
}

// 存款(將錢存至測試帳號)
func MoneyIn(traceCode string, account, amount, currency string) ([]byte, entity.MoneyInResponse, bool) {
	var (
		data entity.MoneyInResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/peace/money/in", cq9GameHost) //{Config.Url}/peace/money/in
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["account"] = account
	formData["amount"] = amount
	formData["currency"] = currency

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// 提款(將錢從測試帳號提出)
func MoneyOut(traceCode string, account, amount, currency string) ([]byte, entity.MoneyOutResponse, bool) {
	var (
		data entity.MoneyOutResponse
		isOK bool
	)

	//設置request內容
	url := fmt.Sprintf("%s/peace/money/out", cq9GameHost) //{Config.Url}/peace/money/out
	header := map[string]string{}
	header[reqireheader.Authorization] = cq9AuthToken
	header[reqireheader.ContentType] = reqireheader.FormUrlEncode
	formData := map[string]string{}
	formData["account"] = account
	formData["amount"] = amount
	formData["currency"] = currency

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Post, url, header, formData)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

//---------後台爬蟲----------

// 後台登入
func BackendLogin(traceCode, account, password string) ([]byte, entity.BackendLoginResponse, bool) {
	var (
		baseResponse entity.CQ9BaseHttpResponse
		data         entity.BackendLoginResponse
		isOK         bool
	)

	//設置request內容
	url := "https://gsos.996688.co/api/login"
	header := map[string]string{}
	header[reqireheader.ContentType] = reqireheader.Json
	bodyModel := entity.BackendLoginRequest{
		Account:  account,
		Password: password,
	}

	body := serializer.JsonMarshal(traceCode, bodyModel)
	//JsonMarshal失敗返回錯誤
	if body == nil {
		return nil, data, false
	}

	//調用API
	resp, isOK := mhttp.CallThirdApiByBody(traceCode, httpmethod.Post, url, header, body)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response第一層解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &baseResponse)
	//CQ9 response第一層解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	//如果CQ9返回的第一層Status.Code!=0(成功),返回錯誤
	if baseResponse.Status.Code != string(errorcode.Success) {
		return resp, data, false
	}

	//CQ9 response第二層解析
	isOK = serializer.JsonUnMarshal(traceCode, baseResponse.Data, &data)
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}

// 後台報表
func BackendReport(traceCode, token, fromDate, toDate string) ([]byte, entity.BackendReportResponse, bool) {
	var (
		data entity.BackendReportResponse
		isOK bool
	)

	//設置request內容
	url := "https://gsos.996688.co/api/games/vendor/analysis/probability/report"
	header := map[string]string{}
	header[reqireheader.ContentType] = reqireheader.Json
	header[reqireheader.Authorization] = fmt.Sprintf(reqireheader.JwtTokenFormat, token)

	fromTime, isOK := mtime.ParseToLocal(traceCode, mtime.BackendReportTimeFormat, fromDate)
	//parse time失敗返回null
	if !isOK {
		return nil, data, false
	}
	toTime, isOK := mtime.ParseToLocal(traceCode, mtime.BackendReportTimeFormat, toDate)
	//parse time失敗返回null
	if !isOK {
		return nil, data, false
	}

	//設置轉query的結構
	reportQuery := entity.BackendReportQuery{
		Currency: currency.CNY, //CNY
		GameCode: "",
		GameTeam: cq9TeamID, //CC
		GameType: "slot",
		GroupBy:  "day",
		Language: "zh-cn",
		Page:     "1",
		PageSize: "15",
		FromDate: mtime.TimeString(fromTime, mtime.ApiTimeFormat),
		ToDate:   mtime.TimeString(toTime, mtime.ApiTimeFormat),
	}
	queryString, isOK := mhttp.ToQueryString(traceCode, reportQuery)
	//轉成query string失敗返回null
	if !isOK {
		return nil, data, false
	}

	//調用API
	resp, isOK := mhttp.CallThirdApi(traceCode, httpmethod.Get, fmt.Sprintf("%s?%s", url, queryString), header, nil)
	//調用失敗返回null
	if !isOK {
		return resp, data, false
	}

	//CQ9 response解析
	isOK = serializer.JsonUnMarshal(traceCode, resp, &data)
	//CQ9 response解析失敗返回null
	if !isOK {
		return resp, data, false
	}

	return resp, data, true
}
