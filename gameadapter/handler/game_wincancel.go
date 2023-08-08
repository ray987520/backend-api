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
////    贏分取消
/////////////////////////////

// 贏分取消Request
type GameWinCancelRequest struct {
	BaseSelfDefine        //自訂headers
	MtCode         string `json:"mtcode" validate:"mtcode"` //交易代碼
}

// 序列化GameWinCancelRequest
func (req *GameWinCancelRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 贏分取消 Response
type GameWinCancelResponse struct {
	Data   GameWinCancelResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus      `json:"status"` //狀態欄
}

// 贏分取消 Response data
type GameWinCancelResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary	Probability Game Win Cancel (贏分取消)
//	@Tags		Probability Game Flow
//	@Accept		x-www-form-urlencoded
//	@Param		mtcode	formData	string	true	"交易代碼"
//	@Success	200		{object}	GameWinCancelResponse
//	@Router		/gamepool/CC/game/win/cancel [post]
//
//	@Security	Bearer
func GameWinCancel(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := GameWinCancelRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.GameWinCancel, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.MtCode = r.FormValue("mtcode")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.GameWinCancel, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.ProbabilityGameWinCancel(traceCode, request.MtCode)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded ProbabilityGameWinCancel failure!")
		return errResp
	}

	//轉換data
	data := GameWinCancelResponse{
		Data: GameWinCancelResponseData(cq9Data.Data),
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
