package handler

import (
	cid "GamePoolApi/backgroundmanager/enum/controller"
	"GamePoolApi/backgroundmanager/enum/errorcode"
	"GamePoolApi/backgroundmanager/enum/respmsg"
	"GamePoolApi/common/database"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	"net/http"
)

//---------注單資訊----------

/////////////////////////////
////    取得注單資訊
/////////////////////////////

// 取得注單資訊Request
type BetSlipInfoRequest struct {
	BaseSelfDefine
	Token string `json:"token"` //Token
}

// 序列化BetSlipInfoRequest
func (req *BetSlipInfoRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 取得注單資訊Response
type BetSlipInfoResponse struct {
	Data   BetSlipInfoResponseData   `json:"Data"`   //資料給予的地方
	Status BackendHttpResponseStatus `json:"Status"` //狀態
}

// 取得注單資訊Response data
type BetSlipInfoResponseData struct {
	PlayerId  string `json:"id"`       //玩家編號
	Account   string `json:"account"`  //玩家帳號  ※字串長度限制36個字元
	PAccount  string `json:"paccount"` //代理帳號 (會由我方依據要求方來源判斷此欄位要不要有值，若無值則不用顯示)
	RoundId   string `json:"roundid"`  //Round ID
	GametType string `json:"gametype"` //遊戲類別
	GameCode  string `json:"gameCode"` //遊戲代碼
}

//	@Summary	取得注單資訊
//	@Tags		注單資訊
//	@Accept		json
//	@Param		token	query		string	true	"Token"
//	@Success	200		{object}	BetSlipInfoResponse
//	@Router		/Api/BetSlipInfo [get]
//
//	@Security	Bearer
func BetSlipInfo(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := BetSlipInfoRequest{}

	//read query string
	request.Token = r.URL.Query().Get("token")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.BetSlipInfo, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := backendError(errorcode.BackendError, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.DetailToken(traceCode, request.Token)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "Send ThirdPartyService(DetailToken) failure!")
		return errResp
	}

	//取DB遊戲詳情
	dbData, isOK := database.GameLogGet(traceCode, cq9Data.Data.RoundId)
	//取DB失敗輸出錯誤
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "DB Platfrom_GameLogGet failure!")
		return errResp
	}

	//設置response
	data := BetSlipInfoResponseData{
		PlayerId:  cq9Data.Data.UserId,
		Account:   cq9Data.Data.Account,
		PAccount:  cq9Data.Data.PAccount,
		GametType: cq9Data.Data.GameType,
		RoundId:   dbData.RoundID,
		GameCode:  dbData.GameCode,
	}
	response := BetSlipInfoResponse{
		Data: data,
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
