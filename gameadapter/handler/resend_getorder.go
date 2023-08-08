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
type GetOrderRequest struct {
	BaseSelfDefine        //自訂headers
	IndexId        string `json:"indexid"` //訂單索引-> roundid:xxx
}

// 序列化GetOrderRequest
func (req *GetOrderRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 查詢注單內容 Response
type GetOrderResponse struct {
	Data   GetOrderResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

// 查詢注單內容 Response data
type GetOrderResponseData struct {
	Order GetOrderResponseDetail `json:"order"` //注單
}

// 查詢注單內容 Response data
type GetOrderResponseDetail struct {
	Account             string           `json:"account"`             //帳號※字串長度限制36個字元
	Balance             float64          `json:"balance"`             //玩家餘額
	Bet                 float64          `json:"bets"`                //本局押分,TOCHECK:文件表格跟範例的json tag不同
	Win                 float64          `json:"wins"`                //win 金額(可為負值),TOCHECK:文件表格跟範例的json tag不同
	WinPc               float64          `json:"winpc"`               //該遊戲玩家從PC贏得的金額
	Jackpots            float64          `json:"jackpots"`            //彩金金額
	JackpotContribution []float64        `json:"jackpotcontribution"` //彩池獎金貢獻值
	JackpotType         string           `json:"jackpottype"`         //彩池獎金類別
	Rake                float64          `json:"rake"`                //抽水金額
	CreateTime          string           `json:"createtime"`          //事件發生時間
	EndRoundTime        string           `json:"endroundtime"`        //DB建立時間
	FinalTime           string           `json:"finaltime"`           //最後異動時間
	GameCode            string           `json:"gamecode"`            //遊戲代號
	GameHall            string           `json:"gamehall"`            //遊戲廠商
	GameToken           string           `json:"gametoken"`           //玩家遊戲token
	GameType            string           `json:"gametype"`            //遊戲類別
	IndexId             string           `json:"indexid"`             //訂單索引
	OrderId             string           `json:"orderid"`             //注單編號,TOCHECK:CQ9 response跟C#代碼有orderid節點,但文件表格沒有範例有
	OwnerId             string           `json:"ownerid"`             //上層代理id
	ParentId            string           `json:"parentid"`            //代理id
	Platform            string           `json:"platform"`            //遊戲平台
	PlayerId            string           `json:"playerid"`            //玩家ID
	RabbitMq            string           `json:"rabbitmq"`            //Rabbit MQ 狀態
	RoundId             string           `json:"roundid"`             //局號
	Detail              []GetOrderDetail `json:"detail"`              //遊戲細節
	SingleRowBet        bool             `json:"singlerowbet"`        //是否為滾輪遊戲
	RoomFee             float64          `json:"roomfee"`             //開房費用
	Status              string           `json:"status"`              //訂單狀態
	EventList           []GetOrderEvent  `json:"eventlist"`           //遊戲行為
	TicketId            string           `json:"ticketid"`            //道具編號
	TicketType          string           `json:"tickettype"`          //免費券類型 1=免費遊戲(獲得一局free game) 2=免費 spin(獲得一次free spin)
	GivenType           string           `json:"giventype"`           //免費券取得類型 1=活動贈送 101=代理贈送 111=寶箱贈送 112=商城購買
	TicketBets          float64          `json:"ticketbets"`          //免費券下注額
	CardWin             float64          `json:"cardwin"`             //派彩加成※最大長度為12位數，及小數點後4位
	UseCard             bool             `json:"usecard"`             //是否派彩加成
}

// GetOrder詳細
type GetOrderDetail struct {
	FreeGame      int64 `json:"freegame"`  //Free game
	LuckyDraw     int64 `json:"luckydraw"` //Lucky draw
	PlBonusatform int64 `json:"bonus"`     //Bonus
}

// GetOrder事件
type GetOrderEvent struct {
	Action    string  `json:"action"`    //行為
	Amount    float64 `json:"amount"`    //贏分金額 ※最大長度為12位數，及小數點後4位
	EventTime string  `json:"eventtime"` //事件發生時間
	MtCode    string  `json:"mtcode"`    //交易代碼
	RecordId  string  `json:"recordid"`  //xxx
}

//	@Summary	Get Order (查詢注單內容)
//	@Tags		Resend mechanism Flow
//	@Param		indexid	query		string	true	"訂單索引-> roundid:xxx"	default(CC123456ab07:cq9)
//	@Success	200		{object}	GetOrderResponse
//	@Router		/gamepool/CC/game/getorder [get]
//
//	@Security	Bearer
func GetOrder(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := GetOrderRequest{}

	//read request query string
	request.IndexId = r.URL.Query().Get("indexid")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.GetOrder, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.GetOrder(traceCode, request.IndexId)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded GetOrder failure!")
		return errResp
	}

	//轉換data
	listDetail := make([]GetOrderDetail, len(cq9Data.Data.Order.Detail))
	for i, v := range cq9Data.Data.Order.Detail {
		listDetail[i] = GetOrderDetail(v)
	}

	listEvent := make([]GetOrderEvent, len(cq9Data.Data.Order.EventList))
	for i, v := range cq9Data.Data.Order.EventList {
		listEvent[i] = GetOrderEvent(v)
	}

	data := GetOrderResponse{
		Data: GetOrderResponseData{
			Order: GetOrderResponseDetail{
				Account:             cq9Data.Data.Order.Account,
				Balance:             cq9Data.Data.Order.Balance,
				Bet:                 cq9Data.Data.Order.Bet,
				Win:                 cq9Data.Data.Order.Win,
				WinPc:               cq9Data.Data.Order.WinPc,
				Jackpots:            cq9Data.Data.Order.Jackpots,
				JackpotContribution: cq9Data.Data.Order.JackpotContribution,
				JackpotType:         cq9Data.Data.Order.JackpotType,
				Rake:                cq9Data.Data.Order.Rake,
				CreateTime:          cq9Data.Data.Order.CreateTime,
				EndRoundTime:        cq9Data.Data.Order.EndRoundTime,
				FinalTime:           cq9Data.Data.Order.FinalTime,
				GameCode:            cq9Data.Data.Order.GameCode,
				GameHall:            cq9Data.Data.Order.GameHall,
				GameToken:           cq9Data.Data.Order.GameToken,
				GameType:            cq9Data.Data.Order.GameType,
				IndexId:             cq9Data.Data.Order.IndexId,
				OrderId:             cq9Data.Data.Order.OrderId,
				OwnerId:             cq9Data.Data.Order.OwnerId,
				ParentId:            cq9Data.Data.Order.ParentId,
				Platform:            cq9Data.Data.Order.Platform,
				PlayerId:            cq9Data.Data.Order.PlayerId,
				RabbitMq:            cq9Data.Data.Order.RabbitMq,
				RoundId:             cq9Data.Data.Order.RoundId,
				Detail:              listDetail,
				SingleRowBet:        cq9Data.Data.Order.SingleRowBet,
				RoomFee:             cq9Data.Data.Order.RoomFee,
				Status:              cq9Data.Data.Order.Status,
				EventList:           listEvent,
				TicketId:            cq9Data.Data.Order.TicketId,
				TicketType:          cq9Data.Data.Order.TicketType,
				GivenType:           cq9Data.Data.Order.GivenType,
				TicketBets:          cq9Data.Data.Order.TicketBets,
				CardWin:             cq9Data.Data.Order.CardWin,
				UseCard:             cq9Data.Data.Order.UseCard,
			},
		},
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
