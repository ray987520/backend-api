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
////    查詢未完成注單
/////////////////////////////

// 查詢未完成注單 Request
type RoundCheckRequest struct {
	BaseSelfDefine        //自訂headers
	FromDate       string `json:"fromdate" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //開始時間
	ToDate         string `json:"todate" validate:"datetime=2006-01-02T15:04:05.999-07:00"`   //結束時間
}

// 序列化RoundCheckRequest
func (req *RoundCheckRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 查詢未完成注單 Response
type RoundCheckResponse struct {
	Data   []string             `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

//	@Summary		Round Check (查詢未完成注單 (RoundID 必須符合格式才能查詢到))
//	@Tags			Resend mechanism Flow
//
//	@Description	1.需於每個小時，定時使用 Round Check 對未完成的回合進行補單、退款行為。 ※每筆資料可以使用Round Detail、Get Order 得知我司狀態 舉例: 當 09:00 時 撈取 07:00~08:00，並每小時以此類推 對區間內的數據做補單、退款。
//	@Description	※若貴司的紀錄為未成單卻在 Get order 撈取到該單據請直接使用注單內之資料將該注單成單
//	@Description	2.補單流程是持續性的機制，須執行至所有未完成的回合成功，流程最少須持續24小時。 舉例: 當 09:00 時 執行 07:00~08:00 補單，並有回合未完成情況下，須至下次補單時繼續執行，直到該回合成功或是24小時後方可停止對該回合的補單。
//	@Description	3.執行補單時，需和當下時間保持一定"時間區間"(依照遊戲自身演繹時間進行推算)，避免對當下進行中遊戲執行補單和退款。
//	@Description	4.補單的datetime請使用補單當下時間
//
//	@Param			fromdate	query		string	true	"開始時間"	default(2023-07-01T00:00:00.000-04:00)
//	@Param			todate		query		string	true	"結束時間"	default(2023-07-24T00:00:00.000-04:00)
//	@Success		200			{object}	RoundCheckResponse
//	@Router			/gamepool/CC/game/roundcheck [get]
//
//	@Security		Bearer
func RoundCheck(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := RoundCheckRequest{}

	//read request query string
	request.FromDate = r.URL.Query().Get("fromdate")
	request.ToDate = r.URL.Query().Get("todate")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.RoundCheck, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.RoundCheck(traceCode, request.FromDate, request.ToDate)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded RoundCheck failure!")
		return errResp
	}

	//轉換data
	data := RoundCheckResponse{
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
