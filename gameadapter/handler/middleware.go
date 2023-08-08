package handler

import (
	"GamePoolApi/common/enum/httpmethod"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/middleware"
	"GamePoolApi/common/enum/reqireheader"
	"GamePoolApi/common/service/mhttp"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/uuid"
	"GamePoolApi/common/service/zaplog"
	"GamePoolApi/gameadapter/cfg"
	"GamePoolApi/gameadapter/enum/errorcode"
	"fmt"
	"net/http"
	"strings"
)

/////////////////////////////
////    Middleware
/////////////////////////////

// IP白名單middleware
func IPWhiteListMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get(reqireheader.CfConnectingIp) //嘗試取第一優先remote ip來源,CF-Connecting-IP
		if ip == "" {
			ip = req.Header.Get(reqireheader.XForwardedFor) //嘗試取次優先remote ip來源,X-Forwarded-For
		}
		if ip == "" {
			ip = req.RemoteAddr //嘗試取最後remote ip來源
		}

		/*
			//ip白名單是否包含ip
			if strings.Contains(cfg.IpWhiteList, ip) {
				response := mhttp.GetHttpResponse(string(errorcode.AuthorizationInvalid), reqTime, traceCode, fmt.Sprintf(authTokenError, token), "")
				data := serializer.JsonMarshal(traceCode, response)
				//response marshal error
				if data == nil {
					data = []byte(responseFormatError)
				}
				err := fmt.Errorf(authTokenError, token)
				zaplog.Errorw(innertrace.MiddlewareError, innertrace.FunctionNode, middleware.AuthMiddleware, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
				writeHttpResponse(w, traceCode, string(errorcode.BadParameters), data)
				return
			}
		*/

		next.ServeHTTP(w, req)
	})
}

// Auth Token驗證middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get(reqireheader.Authorization) //雖然用Bearer Auth,但是C#代碼沒有加Bearer才會通
		reqTime := req.Header.Get(innertrace.RequestTimeNode)
		traceCode := req.Header.Get(innertrace.TraceNode)

		//若token與設置的cq9 auth totken不符,返回認證錯誤
		if cfg.Cq9AuthToken != token {
			data := responseError(errorcode.AuthorizationInvalid, reqTime, traceCode, fmt.Sprintf(authTokenError, token))
			//response marshal error
			if data == nil {
				data = []byte(responseFormatError)
			}
			err := fmt.Errorf(authTokenError, token)
			zaplog.Errorw(innertrace.MiddlewareError, innertrace.FunctionNode, middleware.AuthMiddleware, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
			writeHttpResponse(w, traceCode, string(errorcode.BadParameters), data)
			return
		}

		next.ServeHTTP(w, req)
	})
}

// content-type驗證middleware
func AcceptMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		contentType := req.Header.Get(reqireheader.ContentType)
		reqTime := req.Header.Get(innertrace.RequestTimeNode)
		traceCode := req.Header.Get(innertrace.TraceNode)

		//request method跟content-type不正確的話輸出錯誤
		if req.Method != httpmethod.Get && !strings.Contains(contentType, reqireheader.FormUrlEncode) {
			data := responseError(errorcode.InnerError, reqTime, traceCode, fmt.Sprintf(contentTypeError, req.Method, contentType))
			if data == nil {
				data = []byte(responseFormatError)
			}
			err := fmt.Errorf(genTraceCodeError, traceCode)
			zaplog.Errorw(innertrace.MiddlewareError, innertrace.FunctionNode, middleware.AcceptMiddleware, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, err)
			writeHttpResponse(w, traceCode, string(errorcode.InnerError), data)
			return
		}

		next.ServeHTTP(w, req)
	})
}

// 添加自訂資料middleware,主要有traceCode/requesttime
func SelfHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		traceCode := uuid.Gen(tracer.DefaultTraceId)
		timeZone := mtime.GetTimeZone(8)
		//唯一的traceCode產生失敗就返回異常
		if traceCode == "" {
			//沒有traceCode,所以用default traceCode,查詢時再依賴default traceCode+request time
			traceCode = tracer.DefaultTraceId
			//記錄原始http request,
			logOriginRequest(req, traceCode)
			reqTime := mtime.TimeStringAndFillZero(mtime.LocalNow().In(timeZone), mtime.ApiTimeFormat)
			data := responseError(errorcode.InnerError, reqTime, traceCode, fmt.Sprintf(genTraceCodeError, traceCode))
			if data == nil {
				data = []byte(responseFormatError)
			}
			err := fmt.Errorf(genTraceCodeError, traceCode)
			zaplog.Errorw(innertrace.MiddlewareError, innertrace.FunctionNode, middleware.SelfHeaderMiddleware, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, err)
			writeHttpResponse(w, traceCode, string(errorcode.SomethingWrong), data)
			return
		}
		//記錄原始http request
		logOriginRequest(req, traceCode)
		req.Header.Add(innertrace.TraceNode, traceCode)
		req.Header.Add(innertrace.RequestTimeNode, mtime.TimeStringAndFillZero(mtime.LocalNow().In(timeZone), mtime.ApiTimeFormat))

		next.ServeHTTP(w, req)
	})
}

// 封裝總耗時中間件,輸出http request總耗用時間
func TotalTimeMiddleware(next http.Handler) http.Handler {
	//取開始處理時間
	requestTime := mtime.LocalNow()

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
		traceCode := w.Header().Get(innertrace.TraceNode)
		useTime := mtime.LocalNow().Sub(requestTime)
		zaplog.Infow(innertrace.InfoNode, innertrace.FunctionNode, middleware.TotalTimeMiddleware, innertrace.TraceNode, traceCode, innertrace.TotalTimeNode, useTime.Seconds())
	})
}

/////////////////////////////
////    Middleware function
/////////////////////////////

// 記錄原始http request
func logOriginRequest(req *http.Request, traceCode string) {
	curl, err := mhttp.HttpRequest2Curl(req)
	zaplog.Infow(innertrace.LogOriginRequest, innertrace.FunctionNode, middleware.LogOriginRequest, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("curl", curl, innertrace.ErrorInfoNode, err))
}

// 404 handler
func Handle404(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "404 error\n")
}
