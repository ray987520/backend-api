package handler

import (
	"GamePoolApi/common/database"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	cid "GamePoolApi/gameadapter/enum/controller"
	"GamePoolApi/gameadapter/enum/errorcode"
	"fmt"
	"net/http"
)

//----------Player----------

/////////////////////////////
////    驗證Game Token
/////////////////////////////

// 驗證Game Token Request
type AuthPlayerRequest struct {
	BaseSelfDefine        //自訂headers
	GameToken      string `json:"gametoken" validate:"gt=0"` //玩家遊戲token
}

// 序列化AuthPlayerRequest
func (req *AuthPlayerRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 驗證Game Token Response
type AuthPlayerResponse struct {
	Data   AuthPlayerResponseData `json:"data"`   //資料給予的地方
	Status GaHttpResponseStatus   `json:"status"` //狀態欄
}

// 驗證Game Token Response data
type AuthPlayerResponseData struct {
	Account      string       `json:"account" validate:"acct"` //玩家帳號※字串長度限制36個字元
	Balance      float64      `json:"balance" validate:"amt"`  //玩家餘額
	BetLevel     int          `json:"betlevel"`                //押注限額
	BetThreshold BetThreshold `json:"betthreshold"`            //限紅閥值
	Cobrand      Cobrand      `json:"cobrand"`                 //聯名遊戲資訊
	Currency     string       `json:"currency"`                //玩家幣別
	GameCode     string       `json:"gamecode"`                //遊戲代碼
	GameHall     string       `json:"gamehall"`                //遊戲廠商
	GamePlat     string       `json:"gameplat"`                //遊戲平台
	GameTech     string       `json:"gametech"`                //使用技術
	GameType     string       `json:"gametype"`                //遊戲類型
	Id           string       `json:"id"`                      //玩家id※此值為唯一值，請勿使用 account 替代 id
	IsTestss     bool         `json:"istestss"`                //測試代理
	OwnerId      string       `json:"ownerid"`                 //上層代理id
	ParentId     string       `json:"parentid"`                //代理id
	WebId        int          `json:"webid"`                   //押注限額表代號
}

// 聯名遊戲資訊
type Cobrand struct {
	CreateAt    string      `json:"createat"`    //createat
	Images      []string    `json:"images"`      //圖片資訊
	OwnerId     string      `json:"ownerid"`     //ownerid
	Parentid    string      `json:"parentid"`    //parentid
	Permissions Permissions `json:"permissions"` //聯名遊戲資訊權限
	UpdateAt    string      `json:"updateat"`    //updateat
}

// 聯名遊戲資訊權限
type Permissions struct {
	Basic  bool `json:"basic"`  //獨家設定基本款權限，若為true，需要去讀取image 物件內 category 為 basic 的相關資訊
	Custom bool `json:"custom"` //自定義權限，若為true，需要去讀取image 物件內 category 為 custom 的相關資訊
}

// 限紅閥值
type BetThreshold struct {
	MaxBetType float64   `json:"bettype_maximum"` //單區最大限額 (限押分類遊戲)
	Default    float64   `json:"default"`         //籌碼預設值
	Maximum    float64   `json:"maximum"`         //籌碼最大值
	Minimum    float64   `json:"minimum"`         //籌碼最小值
	MaxRound   float64   `json:"round_maximum"`   //單場總限注額 (彩票類遊戲使用)
	RoomBet    []float64 `json:"roombet"`         //最少為空值，最多會有 6 個元素
}

//	@Summary		Auth (驗證 Game Token，並且回傳此玩家的所有設定資訊)
//	@Tags			Player
//
//	@Description	ownerid = parentid 表示此parentid為總代理 ownerid ≠ parentid 表示此parentid為子代理 請務必紀錄parentid與ownerid，否則將產生對帳的疑慮
//
//	@Accept			x-www-form-urlencoded
//	@Param			gametoken	formData	string	true	"玩家遊戲token"
//	@Success		200			{object}	AuthPlayerResponse
//	@Router			/gamepool/CC/player/auth [post]
//
//	@Security		Bearer
func AuthPlayer(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := AuthPlayerRequest{}

	//read request formdata
	err := r.ParseForm()
	//parse formdata失敗,輸出錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ParseFormError, innertrace.FunctionNode, cid.AuthPlayer, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "parse formdata error!")
		return errResp
	}
	request.GameToken = r.FormValue("gametoken")

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.AuthPlayer, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//驗證request失敗輸出錯誤
	if !validator.IsValidStruct(traceCode, request) {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "bad request data!")
		return errResp
	}

	//轉發 CQ server
	_, cq9Data, isOK := cq9.Auth(traceCode, request.GameToken)
	//轉發失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, "forwarded Auth failure!")
		return errResp
	}

	//add player to db
	resultCode := database.MemberAdd(traceCode, cq9Data.Data.Account, cq9Data.Data.Id, cq9Data.Data.OwnerId, cq9Data.Data.ParentId, cq9Data.Data.BetLevel, cq9Data.Data.Currency, cq9Data.Data.WebId)
	switch resultCode {
	case 0, 1: //繼續執行,0=成功,1=會員已存在
		break
	default:
		errResp := responseError(errorcode.InnerError, requestTime, traceCode, fmt.Sprintf("db API_MemberAdd error! ResultCode:%d", resultCode))
		return errResp
	}

	//轉換data,data有巢狀結構沒辦法直接轉,需要一層一層指定
	data := AuthPlayerResponse{
		Data: AuthPlayerResponseData{
			Account:      cq9Data.Data.Account,
			Balance:      cq9Data.Data.Balance,
			BetLevel:     cq9Data.Data.BetLevel,
			BetThreshold: BetThreshold(cq9Data.Data.BetThreshold),
			Cobrand: Cobrand{
				CreateAt:    cq9Data.Data.Cobrand.CreateAt,
				Images:      cq9Data.Data.Cobrand.Images,
				OwnerId:     cq9Data.Data.OwnerId,
				Parentid:    cq9Data.Data.ParentId,
				Permissions: Permissions(cq9Data.Data.Cobrand.Permissions),
				UpdateAt:    cq9Data.Data.Cobrand.UpdateAt,
			},
			Currency: cq9Data.Data.Currency,
			GameCode: cq9Data.Data.GameCode,
			GameHall: cq9Data.Data.GameHall,
			GamePlat: cq9Data.Data.GamePlat,
			GameTech: cq9Data.Data.GameTech,
			GameType: cq9Data.Data.GameType,
			Id:       cq9Data.Data.Id,
			IsTestss: cq9Data.Data.IsTestss,
			OwnerId:  cq9Data.Data.OwnerId,
			ParentId: cq9Data.Data.ParentId,
			WebId:    cq9Data.Data.WebId,
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
