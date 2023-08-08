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

//---------Resend mechanism Flow----------

/////////////////////////////
////    重新激活相對應 RoundID 的 gametoken
/////////////////////////////

// 重新激活相對應 RoundID 的 gametoken Request
type RecoverGametokenByRoundIDRequest struct {
	BaseSelfDefine        //自訂headers
	IndexId        string `json:"indexid"` //局號
}

// 序列化RecoverGametokenByRoundIDRequest
func (req *RecoverGametokenByRoundIDRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 重新激活相對應 RoundID 的 gametoken Response
type RecoverGametokenByRoundIDResponse struct {
	Data   string               `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

//	@Summary	Recover Gametoken by RoundID (當重送 End 時，遇到 error code 9 Invalid GameToken.，即可使用這個API重新激活相對應)
//	@Tags		Resend mechanism Flow
//	@Accept		x-www-form-urlencoded
//	@Param		indexid	formData	string	true	"訂單索引-> roundid:xxx"	default(CC123456ab05:cq9)
//	@Success	200		{object}	RecoverGametokenByRoundIDResponse
//	@Router		/gamepool/CC/game/generateoneroundcache [post]
//
//	@Security	Bearer
func RecoverGametokenByRoundIDRange(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := RecoverGametokenByRoundIDRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.RecoverGametokenByRoundIDRange, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.IndexId = r.FormValue("indexid")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.RecoverGametokenByRoundIDRange, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.RecoverGametokenByRoundID(traceCode, request.IndexId)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded RecoverGametokenByRoundID failure!")
		return errResp
	}

	//轉換data
	data := RecoverGametokenByRoundIDResponse{
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
