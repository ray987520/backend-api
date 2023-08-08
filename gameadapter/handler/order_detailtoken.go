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

//---------Order Detail(Game Result)----------

/////////////////////////////
////    驗證細單Token正確性並取回相關細單資訊
/////////////////////////////

// 驗證細單Token正確性並取回相關細單資訊Request
type DetailTokenRequest struct {
	BaseSelfDefine        //自訂headers
	Token          string `json:"token" validate:"gt=0"` //細單連結token
}

// 序列化DetailTokenRequest
func (req *DetailTokenRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 驗證細單Token正確性並取回相關細單資訊 Response
type DetailTokenResponse struct {
	Data   DetailTokenResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus    `json:"status"` //狀態欄
}

// 驗證細單Token正確性並取回相關細單資訊 Response data
type DetailTokenResponseData struct {
	RoundId  string `json:"roundid"`            //Round ID
	Account  string `json:"account"`            //玩家帳號※字串長度限制36個字元
	PAccount string `json:"paccount,omitempty"` //代理帳號 (會由我方依據要求方來源判斷此欄位要不要有值，若無值則不用顯示)
	UserId   string `json:"id"`                 //玩家id
	GameType string `json:"gametype"`           //遊戲類別
}

//	@Summary	Detail Token (驗證細單Token正確性並取回相關細單資訊)
//	@Tags		Order Detail(Game Result)
//	@Accept		x-www-form-urlencoded
//	@Param		token	formData	string	true	"細單連結token"
//	@Success	200		{object}	DetailTokenResponse
//	@Router		/gamepool/CC/game/detailtoken [post]
//
//	@Security	Bearer
func DetailToken(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := DetailTokenRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.DetailToken, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.Token = r.FormValue("token")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.DetailToken, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.DetailToken(traceCode, request.Token)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded DetailToken failure!")
		return errResp
	}

	//轉換data
	data := DetailTokenResponse{
		Data: DetailTokenResponseData(cq9Data.Data),
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
