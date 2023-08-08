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
////    查詢注單內容
/////////////////////////////

// 查詢注單內容Request
type RoundDetailRequest struct {
	BaseSelfDefine        //自訂headers
	IndexId        string `json:"indexid" validate:"gt=0"` //訂單索引-> roundid:xxx
}

// 序列化RoundDetailRequest
func (req *RoundDetailRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 查詢注單內容 Response
type RoundDetailResponse struct {
	Data   []RoundDetailResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus      `json:"status"` //狀態欄
}

// 查詢注單內容 Response data
type RoundDetailResponseData struct {
	GameHall   string           `json:"gamehall"`   //遊戲廠商
	GameType   string           `json:"gametype"`   //遊戲種類
	Platform   string           `json:"platform"`   //遊戲平台
	GameCode   string           `json:"gamecode"`   //遊戲代碼
	Account    string           `json:"account"`    //帳號※字串長度限制36個字元
	OwnerId    string           `json:"ownerid"`    //總代理ID
	ParentId   string           `json:"parentid"`   //代理ID
	PlayerId   string           `json:"playerid"`   //玩家ID
	GameToken  string           `json:"gametoken"`  //遊戲Token
	IndexId    string           `json:"indexid"`    //訂單索引
	Status     string           `json:"status"`     //訂單狀態
	CreateTime string           `json:"createtime"` //事件發生時間
	Event      RoundDetailEvent `json:"event"`      //事件
	RoomFee    float64          `json:"roomfee"`    //開房費用
	TicketId   string           `json:"ticketid"`   //道具編號
}

// 注單事件
type RoundDetailEvent struct {
	Action    string  `json:"action"`    //行為
	MtCode    string  `json:"mtcode"`    //交易代碼
	Amount    float64 `json:"amount"`    //交易量
	EventTime string  `json:"eventtime"` //時間
}

//	@Summary		Round Detail (查詢注單內容)
//	@Tags			Resend mechanism Flow
//
//	@Description	依照查看到的狀態去做相對應的動作
//	@Description	*Bet complete、init，且你的遊戲DB沒有 Win 的資料或相對應的後續 ➔ Refund
//	@Description	*有Bet complete 且 Win init，且你的遊戲DB有相對應的資料 ➔ End
//	@Description	*Rollout complete、init，但你的遊戲DB沒有相對應資料 ➔ refund
//	@Description	*Rollout complete 且 Rollin init，且你的遊戲DB有相對應的資料 ➔ Rollin
//	@Description	*Bet / Rollout failure 不執行任何動作
//	@Description	*Bet / Rollout status 為 refunded 不執行任何動作
//
//	@Param			indexid	query		string	true	"訂單索引-> roundid:xxx"	default(CC12334ab:cq9)
//	@Success		200		{object}	RoundDetailResponse
//	@Router			/gamepool/CC/game/rounddetail [get]
//
//	@Security		Bearer
func RoundDetail(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := RoundDetailRequest{}

	//read request query string
	request.IndexId = r.URL.Query().Get("indexid")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.RoundDetail, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.RoundDetail(traceCode, request.IndexId)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded RoundDetail failure!")
		return errResp
	}

	//轉換data
	list := make([]RoundDetailResponseData, len(cq9Data.Data))
	for i, v := range cq9Data.Data {
		list[i] = RoundDetailResponseData{
			GameHall:   v.GameHall,
			GameType:   v.GameType,
			Platform:   v.Platform,
			GameCode:   v.GameCode,
			Account:    v.Account,
			OwnerId:    v.OwnerId,
			ParentId:   v.ParentId,
			PlayerId:   v.PlayerId,
			GameToken:  v.GameToken,
			IndexId:    request.IndexId,
			Status:     v.Status,
			CreateTime: v.CreateTime,
			Event:      RoundDetailEvent(v.Event),
			RoomFee:    v.RoomFee,
			TicketId:   v.TicketId,
		}
	}

	data := RoundDetailResponse{
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
