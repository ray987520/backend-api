package handler

import (
	"GamePoolApi/backgroundmanager/cfg"
	cid "GamePoolApi/backgroundmanager/enum/controller"
	"GamePoolApi/backgroundmanager/enum/errorcode"
	"GamePoolApi/backgroundmanager/enum/respmsg"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/crypt"
	"GamePoolApi/common/service/mhttp"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	"net/http"
)

//---------登入授權令牌、換發----------

/////////////////////////////
////    變更時區
/////////////////////////////

// 變更時區Request
type ChangeTimeZoneRequest struct {
	BaseSelfDefine
	ChangeTimeZoneRequestBody
}

// 變更時區RequestBody
type ChangeTimeZoneRequestBody struct {
	TimeZone string `json:"TimeSpan" default:"+08:00" validate:"tz"` //時區
}

// 序列化ChangeTimeZoneRequest
func (req *ChangeTimeZoneRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 變更時區Response
type ChangeTimeZoneResponse struct {
	AvatarPath  string                    `json:"AvatarPath,omitempty"`   //頭像路徑
	AccessToken string                    `json:"access_token,omitempty"` //訪問令牌
	Status      BackendHttpResponseStatus `json:"Status"`                 //狀態
}

//	@Summary	變更時區
//	@Tags		登入授權令牌、換發
//	@Accept		json
//	@Param		Body	body		ChangeTimeZoneRequestBody	true	"RequestBody"
//	@Success	200		{object}	ChangeTimeZoneResponse
//	@Router		/Token/ChangeTimeZone [post]
//
//	@Security	Bearer
func ChangeTimeZone(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := ChangeTimeZoneRequest{}

	//read request body
	body, isOK := mhttp.ReadHttpRequestBody(traceCode, r)
	//讀取body失敗返回失敗
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "read request body error!")
		return errResp
	}

	isOK = serializer.JsonUnMarshal(traceCode, body, &request)
	//parse body失敗返回失敗
	if !isOK {
		errResp := backendError(errorcode.BackendError, traceCode, "parse request body error!")
		return errResp
	}

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.ChangeTimeZone, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出失敗
	if !validator.IsValidStruct(traceCode, request) {
		errResp := backendError(errorcode.BackendError, traceCode, "bad request data!")
		return errResp
	}

	//產生JWT token/refresh token,refreshToken沒用到先mark
	accessToken := crypt.JwtCreateAccessToken(traceCode, "1", cfg.BackendAdminAccount, request.TimeZone)
	//refreshToken := fmt.Sprintf("%d-%s", mtime.UtcNow().Add(168*time.Hour).Unix(), uuid.Gen(traceCode))

	//設置response
	data := ChangeTimeZoneResponse{
		AvatarPath:  avatarPath,
		AccessToken: accessToken,
		Status: BackendHttpResponseStatus{
			Code:      string(errorcode.Success),
			Message:   respmsg.Success,
			Timestamp: mtime.UtcNow().Unix(),
			//TraceCode: traceCode,
		},
	}
	byteResponse := serializer.JsonMarshal(traceCode, data)

	return byteResponse
}
