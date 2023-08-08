package handler

import (
	"GamePoolApi/backgroundmanager/cfg"
	cid "GamePoolApi/backgroundmanager/enum/controller"
	"GamePoolApi/backgroundmanager/enum/errorcode"
	"GamePoolApi/common/enum/appmode"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/mhttp"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/validator"
	"GamePoolApi/common/service/zaplog"
	"net/http"
)

//---------報表----------

/////////////////////////////
////    取得後台報表
/////////////////////////////

// 取得後台報表Request
type BackendReportRequest struct {
	BaseSelfDefine
	BackendReportRequestBody
}

// 取得後台報表RequestBody
type BackendReportRequestBody struct {
	Token     string `json:"Token" default:"258EAFA5-E914-47DA-95CA-C5AB0DC85B11" validate:"gt=0"` //Token
	StartTime string `json:"StartTime" default:"2023-07-01T00:00:00.000Z" validate:"utc"`          //StartTime,UTC
	EndTime   string `json:"EndTime" default:"2023-07-02T00:00:00.000Z" validate:"utc"`            //EndTime,UTC
}

// 序列化BackendReportRequest
func (req *BackendReportRequest) ToString() string {
	data := serializer.JsonMarshal(req.TraceCode, req)
	return string(data)
}

// 取得後台報表Response
type BackendReportResponse struct {
	Status MaHttpResponseStatus      `json:"status"` //狀態欄
	Data   BackendReportResponseBody `json:"data"`   //資料
}

// 取得後台報表ResponseBody
type BackendReportResponseBody struct {
	TotalBet     float64               `json:"totalBet"`     //總下注
	TotalWin     float64               `json:"totalWin"`     //總贏分
	TotalJackpot float64               `json:"totalJackpot"` //總彩池
	TotalRake    float64               `json:"totalRake"`    //總抽水
	TotalIncome  float64               `json:"totalIncome"`  //總收入
	TotalRtp     float64               `json:"totalRtp"`     //總RTP
	TotalSize    int                   `json:"totalSize"`    //總筆數
	DetailList   []BackendReportDetail `json:"list"`         //列表
}

// 取得後台報表ResponseBody.DetailList
type BackendReportDetail struct {
	Date     string  `json:"date"`     //日期
	GameTeam string  `json:"gameTeam"` //團隊
	GameCode string  `json:"gameCode"` //遊戲代碼
	GameName string  `json:"gameName"` //遊戲名稱
	Count    int     `json:"count"`    //筆數
	Bet      float64 `json:"bet"`      //下注
	Win      float64 `json:"win"`      //贏分
	Jackpot  int64   `json:"jackpot"`  //彩池
	Rake     int64   `json:"rake"`     //抽水
	Income   float64 `json:"income"`   //收入
	Rtp      float64 `json:"rtp"`      //RTP
}

// @Summary	取得後台報表
// @Tags		報表
// @Accept		json
// @Param		Body	body		BackendReportRequestBody	true	"RequestBody"
// @Success	200		{object}	BackendReportResponse
// @Router		/backend/report [post]
//
// @Security	Bearer
func BackendReport(traceCode, requestTime string, r *http.Request) []byte {
	//#region 讀取request

	request := BackendReportRequest{}

	//read request body
	body, isOK := mhttp.ReadHttpRequestBody(traceCode, r)
	//讀取body失敗返回null
	if !isOK {
		return nil
	}

	isOK = serializer.JsonUnMarshal(traceCode, body, &request)
	//parse body失敗返回null
	if !isOK {
		return nil
	}

	//#endregion

	//log request model
	zaplog.Infow(innertrace.LogRequestModel, innertrace.FunctionNode, cid.BackendReport, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("request", string(serializer.JsonMarshal(traceCode, request))))

	//非正式環境返回null
	if cfg.Mode != appmode.Production {
		return nil
	}

	//驗證request失敗輸出null
	if !validator.IsValidStruct(traceCode, request) {
		return nil
	}

	//如果token不是"258EAFA5-E914-47DA-95CA-C5AB0DC85B11"輸出null
	if request.Token != cfg.BackendToken {
		return nil
	}

	//登入後台
	_, data, isOK := cq9.BackendLogin(traceCode, cfg.BackendLoginAccount, cfg.BackendLoginPassword)
	//登入失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.BackendError, requestTime, traceCode, "web crawler failure, Authorization error.")
		return errResp
	}

	//取後台報表
	_, reportData, isOK := cq9.BackendReport(traceCode, data.Token, request.StartTime, request.EndTime)
	//登入失敗輸出錯誤
	if !isOK {
		errResp := responseError(errorcode.BackendError, requestTime, traceCode, "web crawler failure, get report error.")
		return errResp
	}

	//轉換CQ9data
	respDetails := make([]BackendReportDetail, len(reportData.Data.DetailList))
	for i, v := range reportData.Data.DetailList {
		respDetails[i] = BackendReportDetail(v)
	}

	respData := BackendReportResponse{
		Data: BackendReportResponseBody{
			TotalBet:     reportData.Data.TotalBet,
			TotalWin:     reportData.Data.TotalWin,
			TotalJackpot: reportData.Data.TotalJackpot,
			TotalRake:    reportData.Data.TotalRake,
			TotalIncome:  reportData.Data.TotalIncome,
			TotalRtp:     reportData.Data.TotalRtp,
			TotalSize:    reportData.Data.TotalSize,
			DetailList:   respDetails,
		},
		Status: MaHttpResponseStatus{
			Code:       reportData.Status.Code,
			Message:    reportData.Status.Message,
			DateTime:   reportData.Status.DateTime,
			TraceCode:  reportData.Status.TraceCode,
			Latency:    reportData.Status.Latency,
			WalletType: reportData.Status.WalletType,
		},
	}
	byteResponse := serializer.JsonMarshal(traceCode, respData)

	return byteResponse
}
