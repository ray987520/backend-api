package database

import (
	"GamePoolApi/common/entity"
	"GamePoolApi/common/enum/innertrace"
	sqlid "GamePoolApi/common/enum/sql"
	iface "GamePoolApi/common/interface"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
)

/////////////////////////////
////    SQL DB 存取
/////////////////////////////

const (
	rollOutIdFormat   = "rollOut-%s"                         //rollOut transactionId Format
	rollInIdFormat    = "rollIn-%s"                          //rollIn transactionId Format
	unknowError       = "Unknow Error"                       //unknow error default message
	rowCountError     = "rowsAffected is not match expected" //非預期輸出筆數
	currencyError     = "unknow currency:%s"                 //未知幣別
	emptyErrorMessage = "get no error message"               //取出空錯誤訊息
	dataError         = "data is not match expected"         //not expected data
	dbInfo            = "database info"                      //記錄trace用資料
)

var sqlDb iface.ISqlService

// 初始化,注入sql client
func InitSqlWorker(db iface.ISqlService) bool {
	sqlDb = db
	return true
}

// API新增會員
func MemberAdd(traceCode string, account, id, ownerId, parentId string, betLevel int, currency string, webId int) int {
	result := struct {
		ResultCode int
		MemberID   int64
	}{}
	//set parameters
	sql := `EXEC sp_API_MemberAdd ?,?,?,?,?,?,?,?`
	params := []interface{}{account, id, ownerId, parentId, betLevel, "", currency, webId}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.MemberAdd, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return -1
	}

	return result.ResultCode
}

// 對方ID轉換我方ID + 驗證會員存在與否
func PlatformMemberIDChangeMemberID(traceCode string, platformMemberId string) (memberId int) {
	//set parameters
	sql := `EXEC sp_Change_PlatformMemberIDChangeMemberID ?`
	params := []interface{}{platformMemberId}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.PlatformMemberIDChangeMemberID, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &memberId, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return -1
	}

	return memberId
}

// 遊戲代碼轉換遊戲ID
func GameCodeChangeGameID(traceCode string, gameCode string) (gameId int) {
	//set parameters
	sql := `EXEC sp_Change_GameCodeChangeGameID ?`
	params := []interface{}{gameCode}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.GameCodeChangeGameID, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &gameId, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return -1
	}

	return gameId
}

// 遊戲新增賽果+Log
func AddGameResult(traceCode string, memberId, gameId int, round string, bet, winLose, win int64, payoutTime string) int {
	result := struct {
		ResultCode int
		PayoutTime string
	}{}
	//set parameters
	sql := `EXEC sp_Game_GameResultAdd ?,?,?,?,?,?,?,?,?,?,?`

	params := []interface{}{memberId, gameId, round, bet, winLose, win, payoutTime, payoutTime, "", "", ""}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.AddGameResult, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return -1
	}

	return result.ResultCode
}

// 遊戲詳情
func GameLogGet(traceCode string, roundId string) (entity.GameLogGet, bool) {
	var result entity.GameLogGet
	//set parameters
	sql := `EXEC sp_Platfrom_GameLogGet ?`

	params := []interface{}{roundId}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.GameLogGet, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return result, false
	}

	return result, true
}

// 後台遊戲報表筆數依結束時間
func GameResultGetCountByBetTime(traceCode string, loginId, memberId int, roundId, startDate, endDate string) int {
	result := struct {
		ResultCode int
		Ct         int //返回筆數
	}{}
	//set parameters
	sql := `EXEC sp_Platform_GameResultGetCountByBetTime ?,?,?,?,?`

	params := []interface{}{loginId, memberId, roundId, startDate, endDate}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.GameResultGetCountByBetTime, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return -1
	}

	return result.Ct
}

// 後台遊戲報表依結束時間
func GameResultGetListByBetTime(traceCode string, loginId, memberId int, roundId, startDate, endDate string, skipRow, showRow int, field, orderType string) ([]entity.GameResultGetListByBetTime, bool) {
	var result []entity.GameResultGetListByBetTime
	//set parameters
	sql := `EXEC sp_Platform_GameResultGetListByBetTime ?,?,?,?,?,?,?,?,?`

	params := []interface{}{loginId, memberId, roundId, startDate, endDate, skipRow, showRow, field, orderType}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.GameResultGetListByBetTime, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return result, false
	}

	return result, true
}

// 後台遊戲報表筆數依結束時間
func GameReportGetCountByEndTime(traceCode string, loginId int, startDate, endDate string) int {
	result := struct {
		ResultCode int
		Ct         int //返回筆數
	}{}
	//set parameters
	sql := `EXEC sp_Platform_GameReportGetCountByEndTime ?,?,?`

	params := []interface{}{loginId, startDate, endDate}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.GameReportGetCountByEndTime, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return -1
	}

	return result.Ct
}

// 後台遊戲報表依結束時間
func GameReportGetListByEndTime(traceCode string, loginId int, startDate, endDate string) ([]entity.GameReportGetListByEndTime, bool) {
	var result []entity.GameReportGetListByEndTime
	//set parameters
	sql := `EXEC sp_Platform_GameReportGetListByEndTime ?,?,?`

	params := []interface{}{loginId, startDate, endDate}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.GameReportGetListByEndTime, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return result, false
	}

	return result, true
}

// 會員列表筆數
func MemberGetCount(traceCode string, loginId, memberId int) int {
	result := struct {
		ResultCode int
		Ct         int //返回筆數
	}{}
	//set parameters
	sql := `EXEC sp_Platform_MemberGetCount ?,?`

	params := []interface{}{loginId, memberId}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.MemberGetCount, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return -1
	}

	return result.Ct
}

// 會員列表清單
func MemberGetList(traceCode string, loginId, memberId, skipRow, showRow int, field, orderType string) ([]entity.MemberGetList, bool) {
	var result []entity.MemberGetList
	//set parameters
	sql := `EXEC sp_Platform_MemberGetList ?,?,?,?`

	params := []interface{}{loginId, memberId, skipRow, showRow}

	//log exec sql info
	zaplog.Infow(dbInfo, innertrace.FunctionNode, sqlid.MemberGetList, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("sql", sql, "params", params))

	//call sp
	rowCount := sqlDb.CallSP(traceCode, &result, sql, params...)
	//底層錯誤
	if rowCount == -1 {
		return result, false
	}

	return result, true
}
