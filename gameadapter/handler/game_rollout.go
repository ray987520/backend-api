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
////    個人錢包轉至遊戲錢包
/////////////////////////////

// 個人錢包轉至遊戲錢包Request
type RolloutRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken" validate:"gt=0"`                                  //玩家遊戲token
	UserId         string `json:"id" validate:"gt=0"`                                         //玩家id
	MtCode         string `json:"mtcode" validate:"mtcode"`                                   //交易代碼
	Round          string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Amount         string `json:"amount"`                                                     //取款金額※最大長度為12位數，及小數點後4位
	RolloutTime    string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //取款時間(UTC-4)
	GameCode       string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	TakeAll        string `json:"takeall" validate:"omitempty,boolean"`                       //是否取用全部餘額(default: false, 若為true，可不傳amount欄位)
}

// 序列化RolloutRequest
func (req *RolloutRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 個人錢包轉至遊戲錢包 Response
type RolloutResponse struct {
	Data   RolloutResponseData  `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// 個人錢包轉至遊戲錢包 Response data
type RolloutResponseData struct {
	Amount   float64 `json:"amount"`   //取款金額
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary		Probability Game Rollout (個人錢包轉至遊戲錢包)
//	@Tags			Probability Game Flow
//
//	@Description	當玩家沒有遊戲行為( bet = 0, win = 0 )需要呼叫 refund
//
//	@Accept			x-www-form-urlencoded
//	@Param			gametoken	formData	string	true	"玩家遊戲token"
//	@Param			id			formData	string	true	"玩家id"
//	@Param			mtcode		formData	string	true	"交易代碼"
//	@Param			round		formData	string	true	"遊戲回合編號"
//	@Param			amount		formData	float64	true	"取款金額※最大長度為12位數，及小數點後4位"
//	@Param			datetime	formData	string	true	"取款時間(UTC-4)"
//	@Param			gamecode	formData	string	true	"遊戲代號"
//	@Param			takeall		formData	bool	false	"是否取用全部餘額(default: false, 若為true，可不傳amount欄位)"
//	@Success		200			{object}	RolloutResponse
//	@Router			/gamepool/CC/rollout [post]
//
//	@Security		Bearer
func Rollout(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := RolloutRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.Rollout, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.GameToken = r.FormValue("gametoken")
	request.UserId = r.FormValue("id")
	request.MtCode = r.FormValue("mtcode")
	request.Round = r.FormValue("round")
	request.Amount = r.FormValue("amount")
	request.RolloutTime = r.FormValue("datetime")
	request.GameCode = r.FormValue("gamecode")
	//非必須欄位
	request.TakeAll = r.FormValue("takeall")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.Rollout, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.ProbabilityGameRollout(traceCode, request.GameToken, request.UserId, request.MtCode, request.Round, request.Amount, request.RolloutTime, request.GameCode, request.TakeAll)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded ProbabilityGameRollout failure!")
		return errResp
	}

	//轉換data
	data := RolloutResponse{
		Data: RolloutResponseData(cq9Data.Data),
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
