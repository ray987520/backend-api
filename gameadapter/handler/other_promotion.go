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
////    檢查目前是否有活動列表
/////////////////////////////

// 檢查目前是否有活動列表Request
type PromotionRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken" validate:"gt=0"` //玩家遊戲token
}

// 序列化PromotionRequest
func (req *PromotionRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 檢查目前是否有活動列表 Response ※注意 若無活動時data為null，status.code為"0"
type PromotionResponse struct {
	Data   []PromotionResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus    `json:"status"` //狀態欄
}

// 檢查目前是否有活動列表 Response data
type PromotionResponseData struct {
	Name         string        `json:"name"`     //活動名稱
	PromotionUrl string        `json:"promourl"` //推廣活動網址
	ImageUrl     string        `json:"imageurl"` //推廣活動icon網址
	HasLink      bool          `json:"haslink"`  //是否為小遊戲連結
	PromotionId  string        `json:"promoid"`  //活動id
	Icon         PromotionIcon `json:"icon"`     //圖片資訊
}

// 推廣活動圖片資訊
type PromotionIcon struct {
	Png  string `json:"png"`  //圖片url
	Json string `json:"json"` //活動id
}

//	@Summary	Promotion(檢查目前是否有活動列表)
//	@Tags		Others
//	@Param		gametoken	query		string	true	"玩家遊戲token"
//	@Success	200			{object}	PromotionResponse
//	@Router		/gamepool/CC/game/promo [get]
//
//	@Security	Bearer
func Promotion(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := PromotionRequest{}

	//read request query string
	request.GameToken = r.URL.Query().Get("gametoken")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.Promotion, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.Promotion(traceCode, request.GameToken)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded Promotion failure!")
		return errResp
	}

	//轉換data
	list := make([]PromotionResponseData, len(cq9Data.Data))
	for i, v := range cq9Data.Data {
		list[i] = PromotionResponseData{
			Name:         v.Name,
			PromotionUrl: v.PromotionUrl,
			ImageUrl:     v.ImageUrl,
			HasLink:      v.HasLink,
			PromotionId:  v.PromotionId,
			Icon:         PromotionIcon(v.Icon),
		}
	}

	data := PromotionResponse{
		Data: list,
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
