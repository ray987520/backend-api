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
	"strings"
)

//---------注單----------

/////////////////////////////
////    注單列表
/////////////////////////////

// 注單列表Request
type BetSlipListRequest struct {
	BaseSelfDefine
	MemberId  int    `json:"CustomNewMemberID"`                   //會員編號,C#代碼是MemberID,後台是CustomNewMemberID
	RoundId   string `json:"RoundID"`                             //局號
	StartDate string `json:"sDate" validate:"gt=0"`               //開始時間
	EndDate   string `json:"eDate" validate:"gt=0"`               //結束時間
	SkipRow   int    `json:"Skip" validate:"gt=-1"`               //跳過筆數
	ShowRow   int    `json:"Show" validate:"gt=0"`                //顯示筆數
	Field     string `json:"Field" validate:"gt=0"`               //排序欄位
	OrderType string `json:"OrderType" validate:"oneof=asc desc"` //排序類型 asc(小->大)、desc(大->小)
}

// 序列化BetSlipListRequest
func (req *BetSlipListRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 注單列表Response
type BetSlipListResponse struct {
	Count  int                       `json:"Count,omitempty"` //資料筆數
	Data   []BetSlipListResponseData `json:"Data"`            //資料給予的地方
	Status BackendHttpResponseStatus `json:"Status"`          //狀態
}

// 注單列表Response data
type BetSlipListResponseData struct {
	StatusID         int     `json:"StatusID"`         //狀態碼	0 正常	1 會員不存在	2 交易失敗	3 寫賽果失敗	4 寫log失敗
	RoundID          string  `json:"RoundID"`          //局號
	PlatformMemberID string  `json:"PlatformMemberID"` //對方會員編號
	MemberAccount    string  `json:"MemberAccount"`    //對方會員帳號
	GameCode         string  `json:"GameCode"`         //遊戲代碼
	GameName         string  `json:"GameName"`         //遊戲名稱
	Bet              float64 `json:"Bet"`              //下注
	WinLose          float64 `json:"WinLose"`          //輸贏
	Payout           float64 `json:"Payout"`           //派彩
	BetTime          string  `json:"BetTime"`          //下注時間
	PayoutTime       string  `json:"PayoutTime"`       //派彩時間
	EndTime          string  `json:"EndTime"`          //結算時間
}

//	@Summary	注單列表
//	@Tags		注單
//	@Accept		json
//	@Param		MemberID	query		int		false	"會員編號"
//	@Param		RoundID		query		string	false	"局號"
//	@Param		sDate		query		string	true	"開始時間"						default(2023-07-01 00:00:00)
//	@Param		eDate		query		string	true	"結束時間"						default(2023-07-24 23:59:59)
//	@Param		Skip		query		string	true	"跳過筆數"						default(0)
//	@Param		Show		query		string	true	"顯示筆數"						default(10)
//	@Param		Field		query		string	true	"排序欄位"						default(MemberAccount))
//	@Param		OrderType	query		string	true	"排序類型 asc(小->大)、desc(大->小)"	default(desc)
//	@Success	200			{object}	BetSlipListResponse
//	@Router		/BetSlip/BetSlipList [get]
//
//	@Security	Bearer
func BetSlipList(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := BetSlipListRequest{}

	//read query string
	if r.URL.Query().Get("MemberID") != "" {
		tempInt, isOK := str.Atoi(traceCode, r.URL.Query().Get("MemberID"))
		//轉換MemberID失敗返回錯誤
		if !isOK {
			errResp := backendError(errorcode.BackendError, traceCode, "transfer MemberID error!")
			return errResp
		}
		request.MemberId = tempInt
	}

	request.RoundId = r.URL.Query().Get("RoundID")
	request.StartDate = r.URL.Query().Get("sDate")
	request.EndDate = r.URL.Query().Get("eDate")

	tempInt, isOK := str.Atoi(traceCode, r.URL.Query().Get("Skip"))
	//轉換Skip失敗返回錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "transfer Skip error!")
		return errResp
	}
	request.SkipRow = tempInt

	tempInt, isOK = str.Atoi(traceCode, r.URL.Query().Get("Show"))
	//轉換Show失敗返回錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "transfer Show error!")
		return errResp
	}
	request.ShowRow = tempInt

	request.Field = r.URL.Query().Get("Field")
	request.OrderType = r.URL.Query().Get("OrderType")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.BetSlipList, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

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
	timeZone, isOK := str.Atoi(traceCode, strings.Split(claim.TimeZone, ":")[0]) //jwt token時區是+08:00之類的字串,須轉為int使用
	//轉換timeZone失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "transfer TimeZone error!")
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
	startDate := fmt.Sprintf("%s %s", mtime.TimeString(startTime, mtime.BackendUtcTimeFormat)+".000", claim.TimeZone)

	//先轉出UTC0 time
	endTime, isOK := mtime.ParseTime(traceCode, mtime.BackendUtcTimeFormat, request.EndDate)
	//轉換EndDate失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "transfer EndDate error!")
		return errResp
	}
	//然後轉出timestring加上claim.TimeZone
	endDate := fmt.Sprintf("%s %s", mtime.TimeString(endTime, mtime.BackendUtcTimeFormat)+".999", claim.TimeZone)

	//取後台遊戲報表筆數
	gameResultCount := database.GameResultGetCountByBetTime(traceCode, loginId, request.MemberId, request.RoundId, startDate, endDate)

	//取後台遊戲報表
	gameResult, isOK := database.GameResultGetListByBetTime(traceCode, loginId, request.MemberId, request.RoundId, startDate, endDate, request.SkipRow, request.ShowRow, request.Field, request.OrderType)
	//取後台遊戲報表失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "DB Platform_GameResultGetListByBetTime failure!")
		return errResp
	}

	//設置response
	data := []BetSlipListResponseData{}
	if len(gameResult) > 0 {
		for _, result := range gameResult {
			endTime := ""
			//轉換時間
			dbBetTime, isOK := mtime.ParseToTimeZone(traceCode, mtime.ApiTimeFormat, result.BetTime, timeZone)
			//轉換dbBetTime失敗輸出錯誤
			if !isOK {
				errResp := backendError(errorcode.BackendError, traceCode, "transfer dbBetTime error!")
				return errResp
			}

			dbPayoutTime, isOK := mtime.ParseToTimeZone(traceCode, mtime.ApiTimeFormat, result.PayoutTime, timeZone)
			//轉換dbPayoutTime失敗輸出錯誤
			if !isOK {
				errResp := backendError(errorcode.BackendError, traceCode, "transfer dbPayoutTime error!")
				return errResp
			}

			if result.EndTime != "" {
				dbEndTime, isOK := mtime.ParseToTimeZone(traceCode, mtime.ApiTimeFormat, result.EndTime, timeZone)
				//轉換dbEndTime失敗輸出錯誤
				if !isOK {
					errResp := backendError(errorcode.BackendError, traceCode, "transfer dbEndTime error!")
					return errResp
				}
				endTime = mtime.TimeStringAndFillZero(dbEndTime, mtime.SysTimeFormat)
			}
			//轉換資料
			data = append(data, BetSlipListResponseData{
				StatusID:         result.StatusID,
				RoundID:          result.RoundID,
				PlatformMemberID: result.PlatformMemberID,
				MemberAccount:    result.MemberAccount,
				GameCode:         result.GameCode,
				GameName:         result.GameName,
				Bet:              float64(result.Bet) / 10000.0,
				WinLose:          float64(result.WinLose) / 10000.0,
				Payout:           float64(result.Payout) / 10000.0,
				BetTime:          mtime.TimeStringAndFillZero(dbBetTime, mtime.SysTimeFormat),
				PayoutTime:       mtime.TimeStringAndFillZero(dbPayoutTime, mtime.SysTimeFormat),
				EndTime:          endTime,
			})
		}
	}
	response := BetSlipListResponse{
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
