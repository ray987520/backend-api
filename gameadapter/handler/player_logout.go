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

//----------Player----------

/////////////////////////////
////    登出
/////////////////////////////

// 登出Request
type LogoutPlayerRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken"` //玩家遊戲token
}

// 序列化LogoutPlayerRequest
func (req *LogoutPlayerRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 登出Response
type LogoutPlayerResponse struct {
	Data   interface{}          `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// TOCHECK:gametoken文件是必填,原代碼不是
//
//	@Summary	Logout (登出)
//	@Tags		Player
//	@Accept		x-www-form-urlencoded
//	@Param		gametoken	formData	string	false	"玩家遊戲token"
//	@Success	200			{object}	LogoutPlayerResponse
//	@Router		/gamepool/CC/player/logout [post]
//
//	@Security	Bearer
func LogoutPlayer(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := LogoutPlayerRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.LogoutPlayer, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.GameToken = r.FormValue("gametoken")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.LogoutPlayer, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.Logout(traceCode, request.GameToken)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded Logout failure!")
		return errResp
	}

	//轉換data
	data := LogoutPlayerResponse{
		Data: cq9Data.Data,
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
