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
	"net/http"
)

//---------會員----------

/////////////////////////////
////    取得會員列表
/////////////////////////////

// 取得會員列表Request
type MemberListRequest struct {
	BaseSelfDefine
	MemberId  int    `json:"MemberID"`                            //會員編號
	SkipRow   int    `json:"Skip" validate:"gt=-1"`               //跳過筆數
	ShowRow   int    `json:"Show" validate:"gt=0"`                //顯示筆數
	Field     string `json:"Field" validate:"gt=0"`               //排序欄位
	OrderType string `json:"OrderType" validate:"oneof=asc desc"` //排序類型 asc(小->大)、desc(大->小)
}

// 序列化MemberListRequest
func (req *MemberListRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 取得會員列表Response
type MemberListResponse struct {
	Count  int                       `json:"Count,omitempty"` //資料筆數
	Data   []MemberListResponseData  `json:"Data"`            //資料給予的地方
	Status BackendHttpResponseStatus `json:"Status"`          //狀態
}

// 取得會員列表Response data
type MemberListResponseData struct {
	MemberID      int64   `json:"MemberID"`      //我方會員編號
	MemberAccount string  `json:"MemberAccount"` //我方會員編號
	Balance       float64 `json:"Balance"`       //玩家餘額
	Currency      string  `json:"Currency"`      //玩家幣別
	PoolID        int     `json:"PoolID"`        //使用中池編號
	NewPoolID     int     `json:"NewPoolID"`     //下次使用池編號
	RTP           float64 `json:"RTP"`           //RTP
}

//	@Summary	取得會員列表
//	@Tags		會員
//	@Accept		json
//	@Param		MemberID	query		int		false	"會員編號"
//	@Param		Skip		query		string	true	"跳過筆數"						default(0)
//	@Param		Show		query		string	true	"顯示筆數"						default(10)
//	@Param		Field		query		string	true	"排序欄位"						default(MemberAccount)
//	@Param		OrderType	query		string	true	"排序類型 asc(小->大)、desc(大->小)"	default(desc)
//	@Success	200			{object}	MemberListResponse
//	@Router		/Member/MemberList [get]
//
//	@Security	Bearer
func MemberList(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := MemberListRequest{}

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
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.MemberList, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

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

	//取會員列表筆數
	memberCount := database.MemberGetCount(traceCode, loginId, request.MemberId)

	//取會員列表清單
	memberList, isOK := database.MemberGetList(traceCode, loginId, request.MemberId, request.SkipRow, request.ShowRow, request.Field, request.OrderType)
	//取會員列表清單失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "DB Platform_MemberGetList failure!")
		return errResp
	}

	//設置response
	data := []MemberListResponseData{}
	if len(memberList) > 0 {
		for _, member := range memberList {
			//轉換資料
			data = append(data, MemberListResponseData{
				MemberID:      member.MemberID,
				MemberAccount: member.MemberAccount,
				Currency:      member.Currency,
				PoolID:        member.PoolID,
				NewPoolID:     member.NewPoolID,
				RTP:           member.RTP,
			})
		}
	}
	response := MemberListResponse{
		Data:  data,
		Count: memberCount,
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
