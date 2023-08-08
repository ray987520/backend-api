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
////    將金額從個人錢包轉至遊戲錢包
/////////////////////////////

// 將金額從個人錢包轉至遊戲錢包Request
type TableRolloutRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken" validate:"gt=0"`                                  //玩家遊戲token
	UserId         string `json:"id" validate:"gt=0"`                                         //玩家id
	MtCode         string `json:"mtcode" validate:"mtcode"`                                   //交易代碼
	Round          string `json:"round" validate:"gt=0,max=30"`                               //遊戲回合編號
	Amount         string `json:"amount"`                                                     //取款金額※最大長度為12位數，及小數點後4位
	RolloutTime    string `json:"datetime" validate:"datetime=2006-01-02T15:04:05.999-07:00"` //取款時間(UTC-4)
	GameCode       string `json:"gamecode" validate:"gt=0"`                                   //遊戲代號
	GameRole       string `json:"gamerole" validate:"oneof=banker player"`                    //玩家角色為庄家(banker) or 閒家(player)
	BankerType     string `json:"bankertype" validate:"oneof=pc human"`                       //對戰玩家是否有真人[pc|human] pc：對戰玩家沒有真人 human：對戰玩家有真人 ※此欄位為牌桌遊戲使用，非牌桌遊戲此欄位值為空字串 ※如果玩家不支持上庄，只存在與系统對玩。則bankertype 為 PC
	TakeAll        string `json:"takeall" validate:"omitempty,boolean"`                       //是否取用全部餘額(default: false, 若為true，可不傳amount欄位)
}

// 序列化TableRolloutRequest
func (req *TableRolloutRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 將金額從個人錢包轉至遊戲錢包 Response
type TableRolloutResponse struct {
	Data   TableRolloutResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus     `json:"status"` //狀態欄
}

// 將金額從個人錢包轉至遊戲錢包 Response data
type TableRolloutResponseData struct {
	Amount   float64 `json:"amount"`   //取款金額
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

//	@Summary		Table Rollout (將金額從個人錢包轉至遊戲錢包)
//	@Tags			Table Flow
//	@Description	當玩家沒有遊戲行為( bet = 0, win = 0 )需要呼叫 refund
//	@Accept			x-www-form-urlencoded
//	@Param			gametoken	formData	string	true	"玩家遊戲token"
//	@Param			id			formData	string	true	"玩家id"
//	@Param			mtcode		formData	string	true	"交易代碼"
//	@Param			round		formData	string	true	"遊戲回合編號"
//	@Param			amount		formData	float64	true	"取款金額※最大長度為12位數，及小數點後4位"
//	@Param			datetime	formData	string	true	"取款時間(UTC-4)"
//	@Param			gamecode	formData	string	true	"遊戲代號"
//	@Param			gamerole	formData	string	true	"玩家角色為庄家(banker) or 閒家(player)"
//	@Param			bankertype	formData	string	true	"對戰玩家是否有真人[pc|human] ※此欄位為牌桌遊戲使用，非牌桌遊戲此欄位值為空字串"
//	@Param			takeall		formData	bool	false	"是否取用全部餘額(default: false, 若為true，可不傳amount欄位)"
//	@Success		200			{object}	TableRolloutResponse
//	@Router			/gamepool/CC/table/rollout [post]
//
//	@Security		Bearer
func TableRollout(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := TableRolloutRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.TableRollout, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
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
	request.GameRole = r.FormValue("gamerole")
	request.BankerType = r.FormValue("bankertype")
	//非必須欄位
	request.TakeAll = r.FormValue("takeall")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.TableRollout, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.TableRollout(traceCode, request.GameToken, request.UserId, request.MtCode, request.Round, request.Amount, request.RolloutTime, request.GameCode, request.GameRole, request.BankerType, request.TakeAll)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded TableRollout failure!")
		return errResp
	}

	//轉換data
	data := TableRolloutResponse{
		Data: TableRolloutResponseData(cq9Data.Data),
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
