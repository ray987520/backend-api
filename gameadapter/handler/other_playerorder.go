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

//---------Others----------

/////////////////////////////
////    取得玩家注單網址
/////////////////////////////

// 取得玩家注單網址Request
type PlayerOrderRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken" validate:"gt=0"` //玩家遊戲token
}

// 序列化PlayerOrderRequest
func (req *PlayerOrderRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 取得玩家注單網址 Response
type PlayerOrderResponse struct {
	Data   PlayerOrderResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus    `json:"status"` //狀態欄
}

// 取得玩家注單網址 Response
type PlayerOrderResponseData struct {
	PlayerOrderUrl string `json:"url"` //玩家注單網址
}

//	@Summary	Player Order URL(取得玩家注單網址)
//	@Tags		Others
//	@Param		gametoken	query		string	true	"玩家遊戲token"
//	@Success	200			{object}	PlayerOrderResponse
//	@Router		/gamepool/CC/game/playerorder [get]
//
//	@Security	Bearer
func PlayerOrder(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := PlayerOrderRequest{}

	//read request query string
	request.GameToken = r.URL.Query().Get("gametoken")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.PlayerOrder, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.PlayerOrder(traceCode, request.GameToken)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded PlayerOrder failure!")
		return errResp
	}

	//轉換data
	data := PlayerOrderResponse{
		Data: PlayerOrderResponseData(cq9Data.Data),
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
