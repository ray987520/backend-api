package entity

import (
	"GamePoolApi/common/service/serializer"
	"encoding/json"
)

/////////////////////////////
////    CQ9 API 結構
/////////////////////////////

// CQ9標準HttpResponse結構
type CQ9BaseHttpResponse struct {
	Data   json.RawMessage       `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 序列化CQ9BaseHttpResponse
func (res *CQ9BaseHttpResponse) ToString() string {
	data := serializer.JsonMarshal(res.Status.TraceCode, res)
	return string(data)
}

// CQ9標準HttpResponse.Status
type CQ9HttpResponseStatus struct {
	Code       string `json:"code"`                 //狀態碼
	Message    string `json:"message"`              //狀態訊息
	DateTime   string `json:"datetime"`             //回應時間
	TraceCode  string `json:"traceCode"`            //追蹤碼
	Latency    string `json:"latency,omitempty"`    //latency
	WalletType string `json:"wallettype,omitempty"` //錢包類別transfer=轉帳錢包，single=單一錢包，ce=虛擬幣錢包
}

//---------GamePool Api----------

// 驗證玩家CQ9 Response
type AuthPlayerResponse struct {
	Data   AuthPlayerResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus  `json:"status"` //狀態欄
}

// 驗證玩家CQ9 Response data
type AuthPlayerResponseData struct {
	Account      string       `json:"account" validate:"acct"` //玩家帳號※字串長度限制36個字元
	Balance      float64      `json:"balance" validate:"amt"`  //玩家餘額
	BetLevel     int          `json:"betlevel"`                //押注限額
	BetThreshold BetThreshold `json:"betthreshold"`            //限紅閥值
	Cobrand      Cobrand      `json:"cobrand"`                 //聯名遊戲資訊
	Currency     string       `json:"currency"`                //玩家幣別
	GameCode     string       `json:"gamecode"`                //遊戲代碼
	GameHall     string       `json:"gamehall"`                //遊戲廠商
	GamePlat     string       `json:"gameplat"`                //遊戲平台
	GameTech     string       `json:"gametech"`                //使用技術
	GameType     string       `json:"gametype"`                //遊戲類型
	Id           string       `json:"id"`                      //玩家id※此值為唯一值，請勿使用 account 替代 id
	IsTestss     bool         `json:"istestss"`                //測試代理
	OwnerId      string       `json:"ownerid"`                 //上層代理id
	ParentId     string       `json:"parentid"`                //代理id
	WebId        int          `json:"webid"`                   //押注限額表代號
}

// 聯名遊戲資訊
type Cobrand struct {
	CreateAt    string      `json:"createat"`    //createat
	Images      []string    `json:"images"`      //圖片資訊
	OwnerId     string      `json:"ownerid"`     //ownerid
	Parentid    string      `json:"parentid"`    //parentid
	Permissions Permissions `json:"permissions"` //聯名遊戲資訊權限
	UpdateAt    string      `json:"updateat"`    //updateat
}

// 聯名遊戲資訊權限
type Permissions struct {
	Basic  bool `json:"basic"`  //獨家設定基本款權限，若為true，需要去讀取image 物件內 category 為 basic 的相關資訊
	Custom bool `json:"custom"` //自定義權限，若為true，需要去讀取image 物件內 category 為 custom 的相關資訊
}

// 限紅閥值
type BetThreshold struct {
	MaxBetType float64   `json:"bettype_maximum"` //單區最大限額 (限押分類遊戲)
	Default    float64   `json:"default"`         //籌碼預設值
	Maximum    float64   `json:"maximum"`         //籌碼最大值
	Minimum    float64   `json:"minimum"`         //籌碼最小值
	MaxRound   float64   `json:"round_maximum"`   //單場總限注額 (彩票類遊戲使用)
	RoomBet    []float64 `json:"roombet"`         //最少為空值，最多會有 6 個元素
}

// 取玩家餘額CQ9 Response
type ShowBalanceResponse struct {
	Data   ShowBalanceResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus   `json:"status"` //狀態欄
}

// 取玩家餘額CQ9 Response data
type ShowBalanceResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 遊戲回合下注CQ9 Response
type GameBetResponse struct {
	Data   GameBetResponseData   `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 遊戲回合下注CQ9 Response data
type GameBetResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 遊戲回合贏分CQ9 Response
type GameWinResponse struct {
	Data   GameWinResponseData   `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 遊戲回合贏分CQ9 Response data
type GameWinResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 遊戲回合結算CQ9 Response
type GameEndResponse struct {
	Data   GameEndResponseData   `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 遊戲回合結算CQ9 Response data
type GameEndResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 贏分取消CQ9 Response
type GameWinCancelResponse struct {
	Data   GameWinCancelResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus     `json:"status"` //狀態欄
}

// 贏分取消CQ9 Response data
type GameWinCancelResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 個人錢包轉至遊戲錢包CQ9 Response
type RolloutResponse struct {
	Data   RolloutResponseData   `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 個人錢包轉至遊戲錢包CQ9 Response data
type RolloutResponseData struct {
	Amount   float64 `json:"amount"`   //取款金額
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 遊戲錢包轉至個人錢包CQ9 Response
type RollinResponse struct {
	Data   RollinResponseData    `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 遊戲錢包轉至個人錢包CQ9 Response data
type RollinResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 將金額從個人錢包轉至遊戲錢包CQ9 Response
type TableRolloutResponse struct {
	Data   TableRolloutResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus    `json:"status"` //狀態欄
}

// 將金額從個人錢包轉至遊戲錢包CQ9 Response data
type TableRolloutResponseData struct {
	Amount   float64 `json:"amount"`   //取款金額
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 遊戲錢包轉至個人錢包CQ9 Response
type TableRollinResponse struct {
	Data   TableRollinResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus   `json:"status"` //狀態欄
}

// 遊戲錢包轉至個人錢包CQ9 Response data
type TableRollinResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 下注取消CQ9 Response
type TableGameRefundResponse struct {
	Data   TableGameRefundResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus       `json:"status"` //狀態欄
}

// 下注取消CQ9 Response data
type TableGameRefundResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 驗證細單Token正確性並取回相關細單資訊CQ9 Response
type DetailTokenResponse struct {
	Data   DetailTokenResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus   `json:"status"` //狀態欄
}

// 驗證細單Token正確性並取回相關細單資訊CQ9 Response data
type DetailTokenResponseData struct {
	RoundId  string `json:"roundid"`            //Round ID
	Account  string `json:"account"`            //玩家帳號※字串長度限制36個字元
	PAccount string `json:"paccount,omitempty"` //代理帳號 (會由我方依據要求方來源判斷此欄位要不要有值，若無值則不用顯示)
	UserId   string `json:"id"`                 //玩家id
	GameType string `json:"gametype"`           //遊戲類別
}

// 查詢未完成注單CQ9 Response
type RoundCheckResponse struct {
	Data   []string              `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 查詢注單內容CQ9 Response
type RoundDetailResponse struct {
	Data   []RoundDetailResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus     `json:"status"` //狀態欄
}

// 查詢注單內容CQ9 Response data
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

// 查詢注單內容CQ9 Response
type GetOrderResponse struct {
	Data   GetOrderResponseData  `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 查詢注單內容CQ9 Response data
type GetOrderResponseData struct {
	Order GetOrderResponseDetail `json:"order"` //注單
}

// 查詢注單內容CQ9 Response data
type GetOrderResponseDetail struct {
	Account             string           `json:"account"`             //帳號※字串長度限制36個字元
	Balance             float64          `json:"balance"`             //玩家餘額
	Bet                 float64          `json:"bets"`                //本局押分
	Win                 float64          `json:"wins"`                //win 金額(可為負值)
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
	OrderId             string           `json:"orderid"`             //注單編號
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

// 重新激活相對應時間區間內的 gametoken CQ9 Response
type RecoverGametokenByTimeRangeResponse struct {
	Data   []string              `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 重新激活相對應 RoundID 的 gametoken CQ9 Response
type RecoverGametokenByRoundIDResponse struct {
	Data   string                `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 注單補扣款CQ9 Response
type OrderDebitResponse struct {
	Data   OrderDebitResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus  `json:"status"` //狀態欄
}

// 注單補扣款CQ9 Response data
type OrderDebitResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 注單補款CQ9 Response
type OrderCreditResponse struct {
	Data   OrderCreditResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus   `json:"status"` //狀態欄
}

// 注單補款CQ9 Response data
type OrderCreditResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 派發紅利CQ9 Response
type BonusResponse struct {
	Data   BonusResponseData     `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 派發紅利CQ9 Response data
type BonusResponseData struct {
	Balance  float64 `json:"balance"`  //玩家餘額
	Currency string  `json:"currency"` //玩家幣別
}

// 幣別列表CQ9 Response
type CurrencyListResponse struct {
	Data   []CurrencyListResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus      `json:"status"` //狀態欄
}

// 幣別列表CQ9 Response data
type CurrencyListResponseData struct {
	Currency      string  `json:"currency"`      //支援幣別
	Rate          float64 `json:"rate"`          //匯率
	RecommendRate float64 `json:"recommendRate"` //建議轉換比率
}

// 取得玩家注單網址CQ9 Response
type PlayerOrderResponse struct {
	Data   PlayerOrderResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus   `json:"status"` //狀態欄
}

// 取得玩家注單網址CQ9 Response
type PlayerOrderResponseData struct {
	PlayerOrderUrl string `json:"url"` //玩家注單網址
}

// 檢查目前是否有活動列表CQ9 Response
type PromotionResponse struct {
	Data   []PromotionResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus   `json:"status"` //狀態欄
}

// 檢查目前是否有活動列表CQ9 Response data
type PromotionResponseData struct {
	Name         string        `json:"name"`     //活動名稱
	PromotionUrl string        `json:"promourl"` //推廣活動網址
	ImageUrl     string        `json:"imageurl"` //推廣活動icon網址
	HasLink      bool          `json:"haslink"`  //是否為小遊戲連結
	PromotionId  string        `json:"promoid"`  //活動id
	Icon         PromotionIcon `json:"icon"`     //圖片資訊
}

// 推廣活動圖片資訊
type PromotionIcon struct {
	Png  string `json:"png"`  //圖片url
	Json string `json:"json"` //活動id
}

// 產生活動連結CQ9 Response
type PromotionLinkResponse struct {
	Data   string                `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

//---------輔助工具----------

// 創建gametoken(用於測試遊戲客端網址)CQ9 Response
type CreateGameTokenResponse struct {
	Data   CreateGameTokenResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus       `json:"status"` //狀態欄
}

// 創建gametoken(用於測試遊戲客端網址)CQ9 Response data
type CreateGameTokenResponseData struct {
	GameToken string `json:"gametoken"` //玩家遊戲token
}

// 創建創建detailtoken(遊戲細單客端網址)CQ9 Response
type CreateDetailTokenResponse struct {
	Data   CreateDetailTokenResponseData `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus         `json:"status"` //狀態欄
}

// 創建detailtoken(遊戲細單客端網址)CQ9 Response data
type CreateDetailTokenResponseData struct {
	DetailToken string `json:"detailtoken"` //遊戲細單客端網址
}

// 存款(將錢存至測試帳號)CQ9 Response
type MoneyInResponse struct {
	Data   MoneyInResponseData   `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 存款(將錢存至測試帳號)CQ9 Response data
type MoneyInResponseData struct {
	Before    float64 `json:"before"`    //交易前餘額
	Balance   float64 `json:"balance"`   //交易後餘額
	Currency  string  `json:"currency"`  //幣別
	Tracecode float64 `json:"tracecode"` //追蹤碼
}

// 提款(將錢從測試帳號提出)CQ9 Response
type MoneyOutResponse struct {
	Data   MoneyOutResponseData  `json:"data"`   //資料給予的地方
	Status CQ9HttpResponseStatus `json:"status"` //狀態欄
}

// 提款(將錢從測試帳號提出)CQ9 Response data
type MoneyOutResponseData struct {
	Before    float64 `json:"before"`    //交易前餘額
	Balance   float64 `json:"balance"`   //交易後餘額
	Currency  string  `json:"currency"`  //幣別
	Tracecode float64 `json:"tracecode"` //追蹤碼
}

//---------後台爬蟲----------

// 後台登入Request
type BackendLoginRequest struct {
	Account  string `json:"account"`  //後台登入帳號,CC
	Password string `json:"password"` //後台登入密碼,1234
}

// 後台登入Response
type BackendLoginResponse struct {
	Token string `json:"token"` //令牌
}

// 取得後台報表query string結構,用於產生query string
type BackendReportQuery struct {
	Currency string `url:"currency"`
	GameCode string `url:"gameCode"`
	GameTeam string `url:"gameTeam"`
	GameType string `url:"gameType"`
	GroupBy  string `url:"groupBy"`
	Language string `url:"language"`
	Page     string `url:"page"`
	PageSize string `url:"pageSize"`
	FromDate string `url:"fromDate"`
	ToDate   string `url:"toDate"`
}

// 取得後台報表Response
type BackendReportResponse struct {
	Status CQ9HttpResponseStatus     `json:"status"` //狀態欄
	Data   BackendReportResponseBody `json:"data"`   //資料
}

// 取得後台報表ResponseBody
type BackendReportResponseBody struct {
	TotalBet     float64               `json:"totalBet"`     //總下注
	TotalWin     float64               `json:"totalWin"`     //總贏分
	TotalJackpot float64               `json:"totalJackpot"` //總彩池
	TotalRake    float64               `json:"totalRake"`    //總抽水
	TotalIncome  float64               `json:"totalIncome"`  //總收入
	TotalRtp     float64               `json:"totalRtp"`     //總RTP
	TotalSize    int                   `json:"totalSize"`    //總筆數
	DetailList   []BackendReportDetail `json:"list"`         //列表
}

// 取得後台報表ResponseBody.DetailList
type BackendReportDetail struct {
	Date     string  `json:"date"`     //日期
	GameTeam string  `json:"gameTeam"` //團隊
	GameCode string  `json:"gameCode"` //遊戲代碼
	GameName string  `json:"gameName"` //遊戲名稱
	Count    int     `json:"count"`    //筆數
	Bet      float64 `json:"bet"`      //下注
	Win      float64 `json:"win"`      //贏分
	Jackpot  int64   `json:"jackpot"`  //彩池
	Rake     int64   `json:"rake"`     //抽水
	Income   float64 `json:"income"`   //收入
	Rtp      float64 `json:"rtp"`      //RTP
}
