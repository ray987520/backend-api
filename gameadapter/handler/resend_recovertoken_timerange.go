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
////    重新激活相對應時間區間內的 gametoken
/////////////////////////////

// 重新激活相對應時間區間內的 gametoken Request
type RecoverGametokenByTimeRangeRequest struct {
	BaseSelfDefine        //自訂headers
	FromDate       string `json:"fromdate" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //開始時間
	ToDate         string `json:"todate" validate:"datetime=2006-01-02T15:04:05.999-07:00"`   //結束時間
}

// 序列化RecoverGametokenByTimeRangeRequest
func (req *RecoverGametokenByTimeRangeRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 重新激活相對應時間區間內的 gametoken Response
type RecoverGametokenByTimeRangeResponse struct {
	Data   []string             `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

//	@Summary	Recover Gametoken by time range (當重送end時，遇到 error code 9 Invalid GameToken,即可使用這個 API 重新激活相對應時間區間內的 gametoken)
//	@Tags		Resend mechanism Flow
//	@Accept		x-www-form-urlencoded
//	@Param		fromdate	formData	string	true	"開始時間"				default(2023-07-01T00:00:00.000-04:00)
//	@Param		todate		formData	string	true	"結束時間(需小於現在時間一個小時)"	default(2023-07-24T00:00:00.000-04:00)
//	@Success	200			{object}	RecoverGametokenByTimeRangeResponse
//	@Router		/gamepool/CC/game/generateroundcache [post]
//
//	@Security	Bearer
func RecoverGametokenByTimeRange(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := RecoverGametokenByTimeRangeRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.RecoverGametokenByTimeRange, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.FromDate = r.FormValue("fromdate")
	request.ToDate = r.FormValue("todate")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.RecoverGametokenByTimeRange, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.RecoverGametokenByTimeRange(traceCode, request.FromDate, request.ToDate)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded RecoverGametokenByTimeRange failure!")
		return errResp
	}

	//轉換data
	data := RecoverGametokenByTimeRangeResponse{
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
