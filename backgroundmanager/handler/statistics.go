package handler

import (
	cid "GamePoolApi/backgroundmanager/enum/controller"
	"GamePoolApi/backgroundmanager/enum/errorcode"
	"GamePoolApi/backgroundmanager/enum/respmsg"
	"GamePoolApi/common/database"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/crypt"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/str"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	"fmt"
	"net/http"
)

//---------遊戲統計----------

/////////////////////////////
////    遊戲綜合統計
/////////////////////////////

// 遊戲綜合統計Request
type AllGameDataStatisticalRequest struct {
	BaseSelfDefine
	StartDate string `json:"sDate" validate:"gt=0"` //開始時間
	EndDate   string `json:"eDate" validate:"gt=0"` //結束時間
}

// 序列化AllGameDataStatisticalRequest
func (req *AllGameDataStatisticalRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 遊戲綜合統計Response
type AllGameDataStatisticalResponse struct {
	Count  int                                  `json:"Count,omitempty"` //資料筆數
	Data   []AllGameDataStatisticalResponseData `json:"Data"`            //資料給予的地方
	Status BackendHttpResponseStatus            `json:"Status"`          //狀態
}

// 遊戲綜合統計Response data
type AllGameDataStatisticalResponseData struct {
	GameCode string  `json:"GameCode"` //遊戲代碼
	GameName string  `json:"GameName"` //遊戲名稱
	Bet      float64 `json:"Bet"`      //下注
	WinLose  float64 `json:"WinLose"`  //輸贏
	Payout   float64 `json:"Payout"`   //派彩
}

//	@Summary	遊戲綜合統計
//	@Tags		遊戲統計
//	@Accept		json
//	@Param		sDate	query		string	true	"開始時間"	default(2023-07-01 00:00:00)
//	@Param		eDate	query		string	true	"結束時間"	default(2023-07-24 00:00:00)
//	@Success	200		{object}	AllGameDataStatisticalResponse
//	@Router		/GameReport/AllGameDataStatistical [get]
//
//	@Security	Bearer
func AllGameDataStatistical(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := AllGameDataStatisticalRequest{}

	//read query string
	request.StartDate = r.URL.Query().Get("sDate")
	request.EndDate = r.URL.Query().Get("eDate")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.AllGameDataStatistical, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := backendError(errorcode.BackendError, traceCode, "bad request data!")
		return errResp
	}

	// 從jwt token取得使用者訊息
	claim := crypt.JwtValidAccessToken(traceCode, getAuthorizationFromRequest(r))
	loginId, isOK := str.Atoi(traceCode, claim.LoginID)
	//轉換loginId失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "transfer LoginID error!")
		return errResp
	}

	//時間格式轉換為用戶timeZone的db time string,同C#代碼
	//先轉出UTC0 time
	startTime, isOK := mtime.ParseTime(traceCode, mtime.BackendUtcTimeFormat, request.StartDate)
	//轉換StartDate失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "transfer StartDate error!")
		return errResp
	}
	//然後轉出timestring加上claim.TimeZone,格式2006-01-02 15:04:05.999 -07:00
	startDate := fmt.Sprintf("%s %s", mtime.TimeString(startTime, mtime.SysTimeFormat), claim.TimeZone)

	//先轉出UTC0 time
	endTime, isOK := mtime.ParseTime(traceCode, mtime.BackendUtcTimeFormat, request.EndDate)
	//轉換EndDate失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "transfer EndDate error!")
		return errResp
	}
	//然後轉出timestring加上claim.TimeZone
	endDate := fmt.Sprintf("%s %s", mtime.TimeString(endTime, mtime.SysTimeFormat), claim.TimeZone)

	//取後台遊戲報表筆數
	gameResultCount := database.GameReportGetCountByEndTime(traceCode, loginId, startDate, endDate)

	//取後台遊戲報表
	gameResult, isOK := database.GameReportGetListByEndTime(traceCode, loginId, startDate, endDate)
	//取後台遊戲報表失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "DB Platform_GameReportGetListByEndTime failure!")
		return errResp
	}

	//設置response
	data := []AllGameDataStatisticalResponseData{}
	if len(gameResult) > 0 {
		for _, result := range gameResult {
			//轉換資料
			data = append(data, AllGameDataStatisticalResponseData{
				GameCode: result.GameCode,
				GameName: result.GameName,
				Bet:      float64(result.Bet) / 10000.0,
				WinLose:  float64(result.WinLose) / 10000.0,
				Payout:   float64(result.Payout) / 10000.0,
			})
		}
	}
	response := AllGameDataStatisticalResponse{
		Data:  data,
		Count: gameResultCount,
		Status: BackendHttpResponseStatus{
			Code:      string(errorcode.Success),
			Message:   respmsg.Success,
			Timestamp: mtime.UtcNow().Unix(),
			//TraceCode: traceCode,
		},
	}
	byteResponse := serializer.JsonMarshal(traceCode, response)

	return byteResponse
}
