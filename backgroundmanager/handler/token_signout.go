package handler

import (
	cid "GamePoolApi/backgroundmanager/enum/controller"
	"GamePoolApi/backgroundmanager/enum/errorcode"
	"GamePoolApi/backgroundmanager/enum/respmsg"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"net/http"
)

//---------登入授權令牌、換發----------

/////////////////////////////
////    會員登出
/////////////////////////////

// 會員登出Request
type SignOutRequest struct {
	BaseSelfDefine
	AccessToken string `json:"access_token"` //訪問令牌
}

// 序列化SignOutRequest
func (req *SignOutRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 會員登出Response
type SignOutResponse struct {
	Status BackendHttpResponseStatus `json:"Status"` //狀態
}

//	@Summary	會員登出
//	@Tags		登入授權令牌、換發
//	@Accept		json
//	@Param		access_token	query		string	false	"訪問令牌"
//	@Success	200				{object}	SignOutResponse
//	@Router		/token/SignOut [get]
//
//	@Security	Bearer
func SignOut(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := SignOutRequest{}

	//read query string
	request.AccessToken = r.URL.Query().Get("access_token")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.SignOut, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//設置response
	response := SignOutResponse{
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
