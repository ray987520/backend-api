package handler

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	"GamePoolApi/gameadapter/cfg"
	cid "GamePoolApi/gameadapter/enum/controller"
	"GamePoolApi/gameadapter/enum/errorcode"
	"fmt"
	"net/http"
)

//---------接入輔助工具----------

/////////////////////////////
////    創建detailtoken(遊戲細單客端網址)
/////////////////////////////

// 創建detailtoken(遊戲細單客端網址)Request
type CreateDetailTokenRequest struct {
	BaseSelfDefine        //自訂headers
	RoundId        string `json:"roundid" validate:"gt=0"`  //遊戲回合編號
	GameCode       string `json:"gamecode" validate:"gt=0"` //遊戲代號
	GameType       string `json:"gametype" validate:"gt=0"` //遊戲類型
}

// 序列化CreateDetailTokenRequest
func (req *CreateDetailTokenRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 創建創建detailtoken(遊戲細單客端網址) Response
type CreateDetailTokenResponse struct {
	Data   CreateDetailTokenResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus          `json:"status"` //狀態欄
}

// 創建detailtoken(遊戲細單客端網址) Response data
type CreateDetailTokenResponseData struct {
	DetailToken     string `json:"detailtoken"`     //遊戲細單token
	SystemDetailUrl string `json:"systemdetailurl"` //系統商遊戲細單客端網址
	PlayerDetailUrl string `json:"playerdetailurl"` //玩家端遊戲細單客端網址
}

//	@Summary		創建detailtoken(遊戲細單客端網址)
//	@Tags			接入輔助工具
//	@Description	roundid 、gamecode 與 gametype 需要是正確的才能產出detailtoken
//	@Description	{相對應站點的URL}/peace/detailtoken?roundid=123456789&gamecode=1&gametype=slot
//	@Param			roundid		query		string	true	"遊戲回合編號"	default(CC123456ab07)
//	@Param			gamecode	query		string	true	"遊戲代碼"		default(CC1001)
//	@Param			gametype	query		string	true	"遊戲類型"		default(fish)
//	@Success		200			{object}	CreateDetailTokenResponse
//	@Router			/peace/detailtoken [get]
//
//	@Security		Bearer
func CreateDetailToken(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := CreateDetailTokenRequest{}

	//read request query string
	request.RoundId = r.URL.Query().Get("roundid")
	request.GameCode = r.URL.Query().Get("gamecode")
	request.GameType = r.URL.Query().Get("gametype")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.CreateDetailToken, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.CreateDetailToken(traceCode, request.RoundId, request.GameCode, request.GameType)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded CreateDetailToken failure!")
		return errResp
	}

	//轉換資料並產生客端網址
	data := CreateDetailTokenResponse{
		Data: CreateDetailTokenResponseData{
			DetailToken:     cq9Data.Data.DetailToken,
			SystemDetailUrl: fmt.Sprintf(cfg.SystemDetailUrl, cq9Data.Data.DetailToken),
			PlayerDetailUrl: fmt.Sprintf(cfg.PlayerDetailUrl, cq9Data.Data.DetailToken),
		},
		Status: GaHttpResponseStatus{
			Code:       cq9Data.Status.Code,
			Message:    cq9Data.Status.Message,
			DateTime:   requestTime,
			TraceCode:  cq9Data.Status.TraceCode,
			Latency:    cq9Data.Status.Latency,
			WalletType: cq9Data.Status.WalletType,
		},
	}
	byteResponse := serializer.JsonMarshal(traceCode, data)

	return byteResponse
}
