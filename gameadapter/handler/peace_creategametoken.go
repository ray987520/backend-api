package handler

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	"GamePoolApi/gameadapter/cfg"
	cid "GamePoolApi/gameadapter/enum/controller"
	"GamePoolApi/gameadapter/enum/errorcode"
	"fmt"
	"net/http"
)

//---------接入輔助工具----------

/////////////////////////////
////    創建gametoken(用於測試遊戲客端網址)
/////////////////////////////

// 創建gametoken(用於測試遊戲客端網址)Request
type CreateGameTokenRequest struct {
	BaseSelfDefine        //自訂headers
	Account        string `json:"account" validate:"acct"`                                                                                   //帳號
	GameType       string `json:"gametype" validate:"gt=0"`                                                                                  //遊戲類型
	GameCode       string `json:"gamecode" validate:"gt=0"`                                                                                  //遊戲代號
	Currency       string `json:"currency" validate:"omitempty,oneof=CNY MYR THB RUB JPY KRW IDR USD IDR(K) VND(K) EUR SGD HKD INR MMK VND"` //幣別
}

// 序列化CreateGameTokenRequest
func (req *CreateGameTokenRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 創建gametoken(用於測試遊戲客端網址) Response
type CreateGameTokenResponse struct {
	Data   CreateGameTokenResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus        `json:"status"` //狀態欄
}

// 創建gametoken(用於測試遊戲客端網址) Response data
type CreateGameTokenResponseData struct {
	GameToken     string `json:"gametoken"`     //玩家遊戲token
	GameClientUrl string `json:"gameclienturl"` //遊戲客端網址
}

//	@Summary		創建gametoken(用於測試遊戲客端網址)
//	@Tags			接入輔助工具
//	@Description	隨機玩家產生gametoken，指定gametype、指定gamecode {相對應站點的URL}/peace/gametoken?account=random&gametype=slot&gamecode=3
//	@Description	指定玩家產生gametoken，指定gametype、指定gamecode {相對應站點的URL}/peace/gametoken?account={指定玩家帳號}&gametype=slot&gamecode=3
//	@Description	指定玩家產生gametoken，指定currency、指定gametype、指定gamecode {相對應站點的URL}/peace/gametoken?account={指定玩家帳號}&gametype=slot&gamecode=3&currency=CNY
//	@Param			account		query		string	true	"玩家帳號(隨機玩家:random)"	default(random)
//	@Param			gametype	query		string	true	"遊戲類型(測試:slot)"		default(fish)
//	@Param			gamecode	query		string	true	"遊戲代碼(測試:CC01)"		default(CC1001)
//	@Param			currency	query		string	false	"幣別(測試:CNY)"
//	@Success		200			{object}	CreateGameTokenResponse
//	@Router			/peace/gametoken [get]
//
//	@Security		Bearer
func CreateGameToken(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := CreateGameTokenRequest{}

	//read request query string
	request.Account = r.URL.Query().Get("account")
	request.GameType = r.URL.Query().Get("gametype")
	request.GameCode = r.URL.Query().Get("gamecode")
	request.Currency = r.URL.Query().Get("currency")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.CreateGameToken, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.CreateGameToken(traceCode, request.Account, request.GameType, request.GameCode, request.Currency)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded GameToken failure!")
		return errResp
	}

	//轉換資料並產生客端網址
	data := CreateGameTokenResponse{
		Data: CreateGameTokenResponseData{
			GameToken:     cq9Data.Data.GameToken,
			GameClientUrl: fmt.Sprintf(cfg.GameClientUrl, request.GameCode, cq9Data.Data.GameToken),
		},
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
