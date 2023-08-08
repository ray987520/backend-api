package handler

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/middleware"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/uuid"
	"GamePoolApi/common/service/zaplog"
	"GamePoolApi/gameadapter/cfg"
	"GamePoolApi/gameadapter/enum/errorcode"
	"fmt"
	"net/http"

	"GamePoolApi/gameadapter/docs"

	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

/////////////////////////////
////    router基礎結構/變數
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
	gamePoolApiPath     = "/gamepool"                                           //gamepool api uri path
	peaceApiPath        = "/peace"                                              //接入輔助工具 api uri path
	pprofPath           = "/debug"                                              //pprof uri path
	logResponse         = "log response header"                                 //log response header
	badErrorCode        = "no error code"                                       //empty error code
	authTokenError      = "Authorization invalid! Authorization:%s"             //auth token error message
	contentTypeError    = "Content-Type invalid! HttpMethod:%s Content-Type:%s" //content-type error message
)

// api router註冊清單
var (
	gamePoolRoutes []Route
	peaceRoutes    []Route
)

/////////////////////////////
////    註冊controller/middleware
/////////////////////////////

// 初始化,註冊所有api controller/middleware跟api path對應
func initGamePoolRoutes() {
	//Player
	registerGamePool("POST", "/CC/player/auth", ForwardHandler(AuthPlayer), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("GET", "/CC/player/balance", ForwardHandler(ShowBalance), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/player/logout", ForwardHandler(LogoutPlayer), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)

	//Probability Game Flow
	registerGamePool("POST", "/CC/game/bet", ForwardHandler(GameBet), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/game/win", ForwardHandler(GameWin), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/game/end", ForwardHandler(GameEnd), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/game/win/cancel", ForwardHandler(GameWinCancel), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/rollout", ForwardHandler(Rollout), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/rollin", ForwardHandler(Rollin), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)

	//Table Flow
	registerGamePool("POST", "/CC/table/rollout", ForwardHandler(TableRollout), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/v2/CC/table/rollin", ForwardHandler(TableRollin), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/game/bet/refund", ForwardHandler(TableGameRefund), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)

	//Order Detail(Game Result)
	registerGamePool("POST", "/CC/game/detailtoken", ForwardHandler(DetailToken), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)

	//Resend mechanism Flow
	registerGamePool("GET", "/CC/game/roundcheck", ForwardHandler(RoundCheck), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("GET", "/CC/game/rounddetail", ForwardHandler(RoundDetail), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("GET", "/CC/game/getorder", ForwardHandler(GetOrder), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/game/generateroundcache", ForwardHandler(RecoverGametokenByTimeRange), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/game/generateoneroundcache", ForwardHandler(RecoverGametokenByRoundIDRange), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)

	//Game Result Excepition Processing
	registerGamePool("POST", "/CC/order/debit", ForwardHandler(OrderDebit), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/order/credit", ForwardHandler(OrderCredit), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/order/bonus", ForwardHandler(Bonus), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)

	//Infos
	registerGamePool("GET", "/CC/game/currency", ForwardHandler(CurrencyList), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)

	//Others
	registerGamePool("GET", "/CC/game/playerorder", ForwardHandler(PlayerOrder), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("GET", "/CC/game/promo", ForwardHandler(Promotion), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerGamePool("POST", "/CC/game/promo/link", ForwardHandler(PromotionLink), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
}

// 初始化,註冊所有api controller/middleware跟api path對應
func initPeaceRoutes() {
	//接入輔助工具
	registerPeace("GET", "/gametoken", ForwardHandler(CreateGameToken), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerPeace("GET", "/detailtoken", ForwardHandler(CreateDetailToken), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerPeace("POST", "/money/in", ForwardHandler(MoneyIn), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
	registerPeace("POST", "/money/out", ForwardHandler(MoneyOut), TotalTimeMiddleware, AuthMiddleware, AcceptMiddleware)
}

/////////////////////////////
////    router導向主體
/////////////////////////////

// 使用mux Router,分不同前路徑規則劃分為swagger|api,使用不同middleware
func NewRouter(version string) http.Handler {
	//init gamepool api routes data
	initGamePoolRoutes()
	//init 接入輔助工具 api routes data
	initPeaceRoutes()

	//init main router
	r := mux.NewRouter()

	//設定swagger ui畫面
	docs.SwaggerInfo.Title = fmt.Sprintf("VA勝利方舟 外部遊戲伺服器對接文件 Env:%s", cfg.Mode)
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

	//game pool api走subrouter加上middleware
	gamePoolApiRouter := r.PathPrefix(gamePoolApiPath).Subrouter()
	//404處理
	gamePoolApiRouter.HandleFunc("/", Handle404)
	for _, route := range gamePoolRoutes {
		apiSubRouter := gamePoolApiRouter.Methods(route.Method).Path(route.Pattern).Subrouter()
		apiSubRouter.NewRoute().Handler(route.Handler)
		apiSubRouter.Use(route.Middlewares...)
	}

	//接入輔助工具api走subrouter加上middleware
	peaceApiRouter := r.PathPrefix(peaceApiPath).Subrouter()
	for _, route := range peaceRoutes {
		apiSubRouter := peaceApiRouter.Methods(route.Method).Path(route.Pattern).Subrouter()
		apiSubRouter.NewRoute().Handler(route.Handler)
		apiSubRouter.Use(route.Middlewares...)
	}

	//cors
	handler := cors.AllowAll().Handler(r)

	return handler
}

// game pool routes註冊url對controller映射及controller前端middlewares,middlewares依添加順序運作,ex:register("get","/a",a,b,c)運作順序為b=>c=>a=>c=>b
func registerGamePool(method, pattern string, handler ForwardHandler, middlewares ...mux.MiddlewareFunc) {
	gamePoolRoutes = append(gamePoolRoutes, Route{method, pattern, handler, middlewares})
}

// peace routes註冊url對controller映射及controller前端middlewares,middlewares依添加順序運作,ex:register("get","/a",a,b,c)運作順序為b=>c=>a=>c=>b
func registerPeace(method, pattern string, handler ForwardHandler, middlewares ...mux.MiddlewareFunc) {
	peaceRoutes = append(peaceRoutes, Route{method, pattern, handler, middlewares})
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
	logOriginRequest(r, traceCode)

	//controller處理返回[]byte
	data := f(traceCode, reqTime, r)

	//輸出結果
	writeHttpResponse(w, traceCode, "", data)
}

//#endregion
