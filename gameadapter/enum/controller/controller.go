package cid

/////////////////////////////
////    API Controller Function Name
/////////////////////////////

type ControllerId string

// controller function,用於zaplog分類
const (
	AuthPlayer                     ControllerId = "AuthPlayer"
	ShowBalance                    ControllerId = "ShowBalance"
	LogoutPlayer                   ControllerId = "LogoutPlayer"
	GameBet                        ControllerId = "GameBet"
	GameWin                        ControllerId = "GameWin"
	GameEnd                        ControllerId = "GameEnd"
	GameWinCancel                  ControllerId = "GameWinCancel"
	Rollout                        ControllerId = "Rollout"
	Rollin                         ControllerId = "Rollin"
	TableRollout                   ControllerId = "TableRollout"
	TableRollin                    ControllerId = "TableRollin"
	TableGameRefund                ControllerId = "TableGameRefund"
	DetailToken                    ControllerId = "DetailToken"
	RoundCheck                     ControllerId = "RoundCheck"
	RoundDetail                    ControllerId = "RoundDetail"
	GetOrder                       ControllerId = "GetOrder"
	RecoverGametokenByTimeRange    ControllerId = "RecoverGametokenByTimeRange"
	RecoverGametokenByRoundIDRange ControllerId = "RecoverGametokenByRoundIDRange"
	OrderDebit                     ControllerId = "OrderDebit"
	OrderCredit                    ControllerId = "OrderCredit"
	Bonus                          ControllerId = "Bonus"
	CurrencyList                   ControllerId = "CurrencyList"
	PlayerOrder                    ControllerId = "PlayerOrder"
	Promotion                      ControllerId = "Promotion"
	PromotionLink                  ControllerId = "PromotionLink"
	CreateGameToken                ControllerId = "CreateGameToken"
	CreateDetailToken              ControllerId = "CreateDetailToken"
	MoneyIn                        ControllerId = "MoneyIn"
	MoneyOut                       ControllerId = "MoneyOut"
)
