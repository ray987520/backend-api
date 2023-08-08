package main

import (
	"GamePoolApi/common/database"
	"GamePoolApi/common/service/cq9"
	"GamePoolApi/common/service/dbo"
	"GamePoolApi/common/service/mconfig"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"GamePoolApi/gameadapter/cfg"
	"GamePoolApi/gameadapter/handler"
	"net/http"
)

const (
	version = "v1.1.2"
)

// 設定swagger Bearer驗證,在swagger ui會顯示一個驗證按鈕輸入Authorization
//
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//
//	@description				輸入Bearer {jwtToken}
func main() {
	defer tracer.PanicTrace(tracer.DefaultTraceId)

	//初始化底層服務
	initBaseServices()

	//初始化api router,然後聆聽
	routers := handler.NewRouter(version)
	http.ListenAndServe(cfg.ListenPort, routers)
}

// 初始化底層服務,原本散落在init()跟變數裡面,不好掌控且可能發生底層沒初始化就先被呼叫
func initBaseServices() {

	///////////////////////
	////   前3步驟會被其他引用,不要動順序
	///////////////////////

	//1.初始化log底層,包含zap
	zaplog.InitZaplog(cfg.ConfigFileName)

	//2.初始化封裝的viper
	mconfig.InitConfigManager(cfg.ConfigFileName)

	//3.初始化GamePoolApiConfig
	appConfig := cfg.InitGamePoolApiConfig()

	//初始化封裝的時間服務
	mtime.InitTimeService(cfg.TimeZone)

	//初始化封裝的cq9服務
	cq9.InitCQ9Service(appConfig)

	//初始化sql db服務
	sqldb := dbo.GetSqlDb(appConfig) //先get client套件實體
	database.InitSqlWorker(sqldb)    //注入給database
}
