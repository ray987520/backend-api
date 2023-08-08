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
	"GamePoolApi/common/service/uuid"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	"fmt"
	"net/http"
	"time"
)

const (
	avatarPath = "/Images/Avatar/0.png"
)

//---------登入授權令牌、換發----------

/////////////////////////////
////    會員登入
/////////////////////////////

// 會員登入Request
type SignInRequest struct {
	BaseSelfDefine
	SignInRequestBody
}

// 會員登入RequestBody
type SignInRequestBody struct {
	Account  string `json:"Account" default:"admin" validate:"mbacct"`      //會員帳號
	Password string `json:"Password" default:"chimera@888" validate:"gt=0"` //會員密碼
	TimeZone string `json:"TimeSpan" default:"+08:00" validate:"tz"`        //時區
	Device   string `json:"Device" default:"Win-PC" validate:"gt=0"`        //登入設備
	Browser  string `json:"Browser" default:"chrome" validate:"gt=0"`       //登入瀏覽器
}

// 序列化SignInRequest
func (req *SignInRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 會員登入Response
type SignInResponse struct {
	AvatarPath   string                    `json:"AvatarPath"`    //頭像路徑
	Authority    int                       `json:"Authority"`     //權限
	AccessToken  string                    `json:"access_token"`  //訪問令牌
	RefreshToken string                    `json:"refresh_token"` //刷新令牌
	Status       BackendHttpResponseStatus `json:"Status"`        //狀態
}

//	@Summary	會員登入
//	@Tags		登入授權令牌、換發
//	@Accept		json
//	@Param		Body	body		SignInRequestBody	true	"RequestBody"
//	@Success	200		{object}	SignInResponse
//	@Router		/Token/SignIn [post]
//
//	@Security	Bearer
func SignIn(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := SignInRequest{}

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
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.SignIn, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗,如果account不是"admin",password不是"chimera@888",輸出失敗
	if !validator.IsValidStruct(traceCode, request) || request.Account != cfg.BackendAdminAccount || request.Password != cfg.BackendAdminPassword {
		errResp := backendError(errorcode.BackendError, traceCode, "DB Platform_Login Error!")
		return errResp
	}

	//產生JWT token/refresh token
	accessToken := crypt.JwtCreateAccessToken(traceCode, "1", request.Account, request.TimeZone)
	refreshToken := fmt.Sprintf("%d-%s", mtime.UtcNow().Add(72*time.Hour).Unix(), uuid.Gen(traceCode))

	//設置response
	data := SignInResponse{
		AvatarPath:   avatarPath,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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
