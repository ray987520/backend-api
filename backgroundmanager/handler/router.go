package handler

import (
	"GamePoolApi/backgroundmanager/enum/errorcode"
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/middleware"
	"fmt"
	"net/http"
	"strings"

	"GamePoolApi/backgroundmanager/cfg"
	"GamePoolApi/backgroundmanager/docs"
	"GamePoolApi/common/service/crypt"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/uuid"
	"GamePoolApi/common/service/zaplog"

	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

/////////////////////////////
////    router基礎屬性/結構
/////////////////////////////

// api router結構
type Route struct {
	Method      string
	Pattern     string
	Handler     ForwardHandler
	Middlewares []mux.MiddlewareFunc
}

const (
	responseFormatError = "http response json format error"                     //response format error message
	genTraceCodeError   = "gen traceCode error,traceCode:%s"                    //gen tracecode error message
	swaggerPath         = "/swagger"                                            //swagger uri path
	backendApiPath      = "/"                                                   //backend api uri path
	pprofPath           = "/debug"                                              //pprof uri path
	logResponse         = "log response header"                                 //log response header
	badErrorCode        = "no error code"                                       //empty error code
	authTokenError      = "Authorization invalid! Authorization:%s"             //auth token error message
	contentTypeError    = "Content-Type invalid! HttpMethod:%s Content-Type:%s" //content-type error message
)

// api router註冊清單
var (
	routes []Route
)

/////////////////////////////
////    註冊controller/middleware
/////////////////////////////

// 初始化,註冊所有api controller/middleware跟api path對應
func initRoutes() {
	//----------報表----------
	register("POST", "/backend/report", ForwardHandler(BackendReport), TotalTimeMiddleware, SelfHeaderMiddleware, AcceptMiddleware)

	//---------登入授權令牌、換發----------
	register("POST", "/token/signin", ForwardHandler(SignIn), TotalTimeMiddleware, SelfHeaderMiddleware, AcceptMiddleware)
	register("POST", "/token/changetimezone", ForwardHandler(ChangeTimeZone), TotalTimeMiddleware, SelfHeaderMiddleware, AuthMiddleware, AcceptMiddleware)
	register("GET", "/token/signout", ForwardHandler(SignOut), TotalTimeMiddleware, SelfHeaderMiddleware, AcceptMiddleware)

	//---------注單資訊----------
	register("GET", "/api/betslipinfo", ForwardHandler(BetSlipInfo), TotalTimeMiddleware, SelfHeaderMiddleware, AcceptMiddleware)

	//---------注單----------
	register("GET", "/betslip/betsliplist", ForwardHandler(BetSlipList), TotalTimeMiddleware, SelfHeaderMiddleware, AuthMiddleware, AcceptMiddleware)
	register("GET", "/betslip/betslipdetails", ForwardHandler(BetSlipDetails), TotalTimeMiddleware, SelfHeaderMiddleware, AcceptMiddleware)

	//---------遊戲統計----------
	register("GET", "/gamereport/allgamedatastatistical", ForwardHandler(AllGameDataStatistical), TotalTimeMiddleware, SelfHeaderMiddleware, AuthMiddleware, AcceptMiddleware)

	//---------會員----------
	register("GET", "/member/memberlist", ForwardHandler(MemberList), TotalTimeMiddleware, SelfHeaderMiddleware, AuthMiddleware, AcceptMiddleware)
}

/////////////////////////////
////    router導向主體
/////////////////////////////

// 使用mux Router,分不同前路徑規則劃分為swagger|api,使用不同middleware
func NewRouter(version string) http.Handler {

	/////////////////////////////
	////   注意調整router/subrouter順序,backend api沒有統一的path,所以是用/作為subrouter的PathPrefix,弄錯順序可能會都導向某一個subrouter去
	/////////////////////////////

	//init gamepool api routes data
	initRoutes()

	//init main router
	r := mux.NewRouter()

	//設定swagger ui畫面
	docs.SwaggerInfo.Title = fmt.Sprintf("單一錢包 管理後台 Env:%s", cfg.Mode)
	docs.SwaggerInfo.Description = "This is a JWT authentication/authorization sample app"
	docs.SwaggerInfo.Version = version
	//swagger走自己的路徑不用經過middleware,/swagger,default寫法:r.PathPrefix(swaggerPath).Handler(httpSwagger.WrapHandler)
	//部分ui畫面可以自訂的寫法如下,可以控制有沒有swagger外框,插入plugin/uiconfig的JS
	r.PathPrefix(swaggerPath).Handler(httpSwagger.Handler(httpSwagger.Layout(httpSwagger.StandaloneLayout), httpSwagger.DeepLinking(true)))

	//pprof走自己的路徑不用經過middleware,/debug
	pprofRouter := r.PathPrefix(pprofPath).Subrouter()
	pprofRouter.Methods("GET").Path("/pprof").HandlerFunc(pprof.Index)
	pprofRouter.Methods("GET").Path("/allocs").Handler(pprof.Handler("allocs"))
	pprofRouter.Methods("GET").Path("/block").Handler(pprof.Handler("block"))
	pprofRouter.Methods("GET").Path("/cmdline").HandlerFunc(pprof.Cmdline)
	pprofRouter.Methods("GET").Path("/goroutine").Handler(pprof.Handler("goroutine"))
	pprofRouter.Methods("GET").Path("/heap").Handler(pprof.Handler("heap"))
	pprofRouter.Methods("GET").Path("/mutex").Handler(pprof.Handler("mutex"))
	pprofRouter.Methods("GET").Path("/profile").HandlerFunc(pprof.Profile)
	pprofRouter.Methods("GET").Path("/threadcreate").Handler(pprof.Handler("threadcreate"))
	pprofRouter.Methods("GET").Path("/trace").HandlerFunc(pprof.Trace)

	//backend api走subrouter加上middleware
	backendApiRouter := r.PathPrefix(backendApiPath).Subrouter()
	for _, route := range routes {
		apiSubRouter := backendApiRouter.Methods(route.Method).Path(route.Pattern).Subrouter()
		apiSubRouter.NewRoute().Handler(route.Handler)
		apiSubRouter.Use(route.Middlewares...)
	}
	//404處理
	backendApiRouter.HandleFunc("/", Handle404)

	//cors
	handler := cors.AllowAll().Handler(r)

	//忽略url path大小寫敏感
	handler = CaselessMatcher(handler)

	return handler
}

/////////////////////////////
////    router function
/////////////////////////////

// 註冊url對controller映射及controller前端middlewares,middlewares依添加順序運作,ex:register("get","/a",a,b,c)運作順序為b=>c=>a=>c=>b
func register(method, pattern string, handler ForwardHandler, middlewares ...mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middlewares})
}

// 404 handler
func Handle404(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "404 error\n")
}

// Basic Auth
func basicAuth(traceCode, authorization string) bool {
	//authorization格式會是"Basic {base64 string}",所以要分割拿後面的資料解碼取得帳號密碼
	splitToken := strings.Split(authorization, " ")
	//若分割後長度不對返回失敗
	if len(splitToken) != 2 {
		zaplog.Errorw(innertrace.MiddlewareError, innertrace.FunctionNode, middleware.BasicAuth, innertrace.TraceNode, traceCode, innertrace.DataNode, "authorization split len error!")
		return false
	}

	decodeData, isOK := crypt.Base64Decode(traceCode, splitToken[1])
	//若base64解碼失敗返回失敗
	if !isOK {
		zaplog.Errorw(innertrace.MiddlewareError, innertrace.FunctionNode, middleware.BasicAuth, innertrace.TraceNode, traceCode, innertrace.DataNode, "authorization Base64Decode error!")
		return false
	}

	//解碼後內容應該為"{account}:{password}",所以要分割使用
	result := strings.Split(string(decodeData), ":")
	//若分割後長度不對或帳號密碼不對返回失敗
	if len(result) != 2 || result[0] != cfg.BackendAdminAccount || result[1] != cfg.BackendAdminPassword {
		zaplog.Errorw(innertrace.MiddlewareError, innertrace.FunctionNode, middleware.BasicAuth, innertrace.TraceNode, traceCode, innertrace.DataNode, "authorization decodeData error!")
		return false
	}

	return true
}

// 忽略url大小寫,因為C# api route是不區分大小寫的,而golang的route match會發生在middleware之前,所以在router設定完後包一層將url路徑轉成小寫
func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//將url路徑轉成小寫
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

//#region 封裝handler,包含TraceCode/RequestTime

// 實作HandlerFunc介面,原本controller要改成對應型態,ex:CreateGameToken(traceCode,requestTime string,r *http.Request)
type ForwardHandler func(string, string, *http.Request) []byte

// 實作HandlerFunc介面的ServeHTTP
func (f ForwardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	traceCode := uuid.Gen(tracer.DefaultTraceId)
	timeZone := mtime.GetTimeZone(8)
	reqTime := mtime.TimeStringAndFillZero(mtime.LocalNow().In(timeZone), mtime.ApiTimeFormat)

	//唯一的traceCode產生失敗就返回異常
	if traceCode == "" {
		//沒有traceCode,所以用default traceCode,查詢時再依賴default traceCode+request time
		traceCode = tracer.DefaultTraceId
		//記錄原始http request,
		logOriginRequest(r, traceCode)
		data := responseError(errorcode.BackendError, reqTime, traceCode, fmt.Sprintf(genTraceCodeError, traceCode))
		if data == nil {
			data = []byte(responseFormatError)
		}
		err := fmt.Errorf(genTraceCodeError, traceCode)
		zaplog.Errorw(innertrace.MiddlewareError, innertrace.FunctionNode, middleware.SelfHeaderMiddleware, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, err)
		writeHttpResponse(w, traceCode, string(errorcode.BackendError), data)
		return
	}

	//記錄原始http request
	logOriginRequest(r, traceCode)

	//controller處理返回[]byte
	data := f(traceCode, reqTime, r)

	//輸出結果
	writeHttpResponse(w, traceCode, "", data)
}

//#endregion
