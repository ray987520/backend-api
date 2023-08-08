package cfg

import (
	"GamePoolApi/common/service/mconfig"
	"time"
)

/////////////////////////////
////    後台API config
/////////////////////////////

type appConfig struct{}

const (
	ConfigFileName       = "backendconfig"                        //config檔名,不要用空字串或config避免viper預設讀config.xxx
	BackendLoginAccount  = "CC"                                   //後台爬蟲登入帳號
	BackendLoginPassword = "1234"                                 //後台爬蟲登入密碼
	BackendToken         = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11" //後台爬蟲登入token
	BackendAdminAccount  = "admin"                                //後台管理者帳號
	BackendAdminPassword = "chimera@888"                          //後台管理者密碼
)

var (
	ListenPort       string        //api聆聽port,如":8080"
	TimeZone         int           //api本地時區,如-4
	SqlConnectString string        //sql master連線字串
	SqlMaxOpenConns  int           //sql最大連線數
	SqlMaxIdleConns  int           //sql最大閒置連線數,沒在用就會關到剩這個數
	SqlMaxIdleSecond time.Duration //sql連線最大閒置秒數,超過閒置時間就會關閉連線
	Cq9GameHost      string        //cq9遊戲api host
	Cq9AuthToken     string        //cq9 auth token
	Cq9GameHall      string        //cq9遊戲廳
	Cq9TeamID        string        //cq9團隊代碼
	Mode             string        //程式環境,Integration:測試 Production:線上
	JwtKey           string        //Jwt密鑰,至少要 16 字元以上
)

// 初始化ma api config
func InitBackendConfig() *appConfig {
	ListenPort = mconfig.GetString("application.listenPort")            //api聆聽port,如":8080"
	TimeZone = mconfig.GetInt("api.timezone")                           //api本地時區,如-4
	SqlConnectString = mconfig.GetString("sql.connectString.master")    //sql master連線字串
	SqlMaxOpenConns = mconfig.GetInt("sql.maxOpenConns")                //sql最大連線數
	SqlMaxIdleConns = mconfig.GetInt("sql.maxIdleConns")                //sql最大閒置連線數,沒在用就會關到剩這個數
	SqlMaxIdleSecond = mconfig.GetDuration("sql.maxIdleSecond")         //sql連線最大閒置秒數,超過閒置時間就會關閉連線
	Cq9GameHost = mconfig.GetString("thirdApi.CQ9Games.Host")           //cq9遊戲api host
	Cq9AuthToken = mconfig.GetString("thirdApi.CQ9Games.Authorization") //cq9 auth token
	Cq9GameHall = mconfig.GetString("thirdApi.CQ9Games.GameHall")       //cq9遊戲廳
	Cq9TeamID = mconfig.GetString("thirdApi.CQ9Games.TeamID")           //cq9團隊代碼
	Mode = mconfig.GetString("application.mode")                        //程式環境,Integration:測試 Production:線上
	JwtKey = mconfig.GetString("jwt.signKey")                           //Jwt密鑰,至少要 16 字元以上
	return &appConfig{}
}

// get cq9遊戲api host
func (config *appConfig) GetCq9GameHost() string {
	return Cq9GameHost
}

// get cq9 auth token
func (config *appConfig) GetCq9AuthToken() string {
	return Cq9AuthToken
}

// get cq9遊戲廳
func (config *appConfig) GetCq9GameHall() string {
	return Cq9GameHall
}

// get cq9團隊代碼
func (config *appConfig) GetCq9TeamID() string {
	return Cq9TeamID
}

// get sql master連線字串
func (config *appConfig) GetSqlConnectString() string {
	return SqlConnectString
}

// get sql最大連線數
func (config *appConfig) GetSqlMaxOpenConns() int {
	return SqlMaxOpenConns
}

// get sql最大閒置連線數
func (config *appConfig) GetSqlMaxIdleConns() int {
	return SqlMaxIdleConns
}

// get sql連線最大閒置秒數
func (config *appConfig) GetSqlMaxIdleSecond() time.Duration {
	return SqlMaxIdleSecond
}
