package handler

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	cid "GamePoolApi/gameadapter/enum/controller"
	"GamePoolApi/gameadapter/enum/errorcode"
	"net/http"
)

//---------Others----------

/////////////////////////////
////    產生活動連結
/////////////////////////////

// 產生活動連結Request
type PromotionLinkRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken" validate:"gt=0"` //玩家遊戲token
	PromotionId    string `json:"promoid" validate:"gt=0"`   //活動id
}

// 序列化PromotionLinkRequest
func (req *PromotionLinkRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 產生活動連結 Response
type PromotionLinkResponse struct {
	Data   string               `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus `json:"status"` //狀態欄
}

//	@Summary	Promotion Link(產生活動連結)
//	@Tags		Others
//	@Accept		x-www-form-urlencoded
//	@Param		gametoken	formData	string	true	"遊戲token"
//	@Param		promoid		formData	string	true	"活動id"
//	@Success	200			{object}	PromotionLinkResponse
//	@Router		/gamepool/CC/game/promo/link [post]
//
//	@Security	Bearer
func PromotionLink(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := PromotionLinkRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.PromotionLink, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.GameToken = r.FormValue("gametoken")
	request.PromotionId = r.FormValue("promoid")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.PromotionLink, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.PromotionLink(traceCode, request.GameToken, request.PromotionId)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded PromotionLink failure!")
		return errResp
	}

	//轉換data
	data := PromotionLinkResponse{
		Data: cq9Data.Data,
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
