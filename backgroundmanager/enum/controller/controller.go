package handler

/////////////////////////////
////    API Controller Function Name
/////////////////////////////

type ControllerId string

// controller function,用於zaplog分類
const (
	BackendReport          ControllerId = "BackendReport"
	SignIn                 ControllerId = "SignIn"
	ChangeTimeZone         ControllerId = "ChangeTimeZone"
	SignOut                ControllerId = "SignOut"
	BetSlipInfo            ControllerId = "BetSlipInfo"
	BetSlipList            ControllerId = "BetSlipList"
	BetSlipDetails         ControllerId = "BetSlipDetails"
	AllGameDataStatistical ControllerId = "AllGameDataStatistical"
	MemberList             ControllerId = "MemberList"
)
