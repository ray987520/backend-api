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
////    下注取消
/////////////////////////////

// 下注取消Request
type TableGameRefundRequest struct {
	BaseSelfDefine        //自訂headers
	MtCode         string `json:"mtcode" validate:"mtcode"` //交易代碼
	UserId         string `json:"userid"`                   //玩家id
}

// 序列化TableGameRefundRequest
func (req *TableGameRefundRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 下注取消 Response
type TableGameRefundResponse struct {
	Data   TableGameRefundResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus        `json:"status"` //狀態欄
}

// 下注取消 Response data
type TableGameRefundResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary	Probability / Table Game Refund (下注取消)
//	@Tags		Table Flow
//	@Accept		x-www-form-urlencoded
//	@Param		mtcode	formData	string	true	"交易代碼"
//	@Param		userid	formData	string	false	"玩家id"
//	@Success	200		{object}	TableGameRefundResponse
//	@Router		/gamepool/CC/game/bet/refund [post]
//
//	@Security	Bearer
func TableGameRefund(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := TableGameRefundRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.TableGameRefund, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.MtCode = r.FormValue("mtcode")
	//非必須欄位
	request.UserId = r.FormValue("userid")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.TableGameRefund, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.TableGameRefund(traceCode, request.MtCode, request.UserId)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded TableGameRefund failure!")
		return errResp
	}

	//轉換data
	data := TableGameRefundResponse{
		Data: TableGameRefundResponseData(cq9Data.Data),
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
