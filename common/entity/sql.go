package entity

/////////////////////////////
////    SQL DB輸出結構
/////////////////////////////

//遊戲詳情model
type GameLogGet struct {
	StatusID         int    `json:"StatusID"`         //狀態碼	0 正常	1 會員不存在	2 交易失敗	3 寫賽果失敗	4 寫log失敗
	OwnerID          string `json:"OwnerID"`          //總代理編號
	ParentID         string `json:"ParentID"`         //代理編號
	PlatformMemberID string `json:"PlatformMemberID"` //對方會員編號
	MemberAccount    string `json:"MemberAccount"`    //對方會員帳號
	GameCode         string `json:"GameCode"`         //遊戲代碼
	GameTypeID       int    `json:"GameTypeID"`       //遊戲類型編號 1:slot 2:fish
	GameName         string `json:"GameName"`         //遊戲名稱
	RoundID          string `json:"RoundID"`          //局號
	Currency         string `json:"Currency"`         //幣別
	Bet              int64  `json:"Bet"`              //下注
	Payout           int64  `json:"Payout"`           //派彩
	GameLog          string `json:"GameLog"`          //遊戲log
	BetTime          string `json:"BetTime"`          //下注時間
	PayoutTime       string `json:"PayoutTime"`       //派彩時間
	EndTime          string `json:"EndTime"`          //結算時間
}

//後台遊戲報表依結束時間model
type GameResultGetListByBetTime struct {
	StatusID         int    `json:"StatusID"`         //狀態碼	0 正常	1 會員不存在	2 交易失敗	3 寫賽果失敗	4 寫log失敗
	RoundID          string `json:"RoundID"`          //局號
	PlatformMemberID string `json:"PlatformMemberID"` //對方會員編號
	MemberAccount    string `json:"MemberAccount"`    //對方會員帳號
	GameCode         string `json:"GameCode"`         //遊戲代碼
	GameName         string `json:"GameName"`         //遊戲名稱
	Bet              int64  `json:"Bet"`              //下注
	WinLose          int64  `json:"WinLose"`          //輸贏
	Payout           int64  `json:"Payout"`           //派彩
	BetTime          string `json:"BetTime"`          //下注時間
	PayoutTime       string `json:"PayoutTime"`       //派彩時間
	EndTime          string `json:"EndTime"`          //結算時間
}

//後台遊戲報表依結束時間
type GameReportGetListByEndTime struct {
	GameCode string `json:"GameCode"` //遊戲代碼
	GameName string `json:"GameName"` //遊戲名稱
	Bet      int64  `json:"Bet"`      //下注
	WinLose  int64  `json:"WinLose"`  //輸贏
	Payout   int64  `json:"Payout"`   //派彩
}

//會員列表清單
type MemberGetList struct {
	MemberID      int64   `json:"MemberID"`      //我方會員編號
	MemberAccount string  `json:"MemberAccount"` //會員帳號
	Currency      string  `json:"Currency"`      //玩家幣別
	PoolID        int     `json:"PoolID"`        //使用中池編號
	NewPoolID     int     `json:"NewPoolID"`     //下次使用池編號
	RTP           float64 `json:"RTP"`           //RTP
}
