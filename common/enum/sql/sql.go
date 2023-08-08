package sql

/////////////////////////////
////    SQL Function Name
/////////////////////////////

type SqlId string

// sql function,用於zaplog分類
const (
	GetExternalErrorMessage        SqlId = "GetExternalErrorMessage"
	MemberAdd                      SqlId = "MemberAdd"
	PlatformMemberIDChangeMemberID SqlId = "PlatformMemberIDChangeMemberID"
	GameCodeChangeGameID           SqlId = "GameCodeChangeGameID"
	AddGameResult                  SqlId = "AddGameResult"
	GameLogGet                     SqlId = "GameLogGet"
	GameResultGetCountByBetTime    SqlId = "GameResultGetCountByBetTime"
	GameResultGetListByBetTime     SqlId = "GameResultGetListByBetTime"
	GameReportGetCountByEndTime    SqlId = "GameReportGetCountByEndTime"
	GameReportGetListByEndTime     SqlId = "GameReportGetListByEndTime"
	MemberGetCount                 SqlId = "MemberGetCount"
	MemberGetList                  SqlId = "MemberGetList"
)
