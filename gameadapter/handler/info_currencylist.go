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

//---------Infos----------

/////////////////////////////
////    幣別列表
/////////////////////////////

// 幣別列表Request
type CurrencyListRequest struct {
	BaseSelfDefine //自訂headers
}

// 序列化CurrencyListRequest
func (req *CurrencyListRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 幣別列表CQ9 Response
type CurrencyListResponse struct {
	Data   []CurrencyListResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus       `json:"status"` //狀態欄
}

// 幣別列表CQ9 Response data
type CurrencyListResponseData struct {
	Currency      string  `json:"currency"`      //支援幣別
	Rate          float64 `json:"rate"`          //匯率
	RecommendRate float64 `json:"recommendRate"` //建議轉換比率
}

//	@Summary	Currency List (幣別列表)
//	@Tags		Infos
//	@Success	200	{object}	CurrencyListResponse
//	@Router		/gamepool/CC/game/currency [get]
//
//	@Security	Bearer
func CurrencyList(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := CurrencyListRequest{}

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.CurrencyList, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.CurrencyList(traceCode)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded CurrencyList failure!")
		return errResp
	}

	//轉換data
	list := make([]CurrencyListResponseData, len(cq9Data.Data))
	for i, v := range cq9Data.Data {
		list[i] = CurrencyListResponseData(v)
	}

	data := CurrencyListResponse{
		Data: list,
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
