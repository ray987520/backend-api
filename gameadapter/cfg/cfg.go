package cfg

import (
	"GamePoolApi/common/service/mconfig"
	"time"
)

/////////////////////////////
////    GA API config
/////////////////////////////

type appConfig struct{}

const (
	ConfigFileName = "gamepoolconfig" //config檔名,不要用空字串或config避免viper預設讀config.xxx
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
	IpWhiteList      string        //IP白名單
	Mode             string        //程式環境,Integration:測試 Production:線上
	GameClientUrl    string        //遊戲客端網址
	SystemDetailUrl  string        //系統商遊戲細單客端網址
	PlayerDetailUrl  string        //玩家端遊戲細單客端網址
)

// 初始化game pool api config
func InitGamePoolApiConfig() *appConfig {
	ListenPort = mconfig.GetString("application.listenPort")            //api聆聽port,如":8080"
	TimeZone = mconfig.GetInt("api.timezone")                           //api本地時區,如-4
	SqlConnectString = mconfig.GetString("sql.connectString.master")    //sql master連線字串
	SqlMaxOpenConns = mconfig.GetInt("sql.maxOpenConns")                //sql最大連線數
	SqlMaxIdleConns = mconfig.GetInt("sql.maxIdleConns")                //sql最大閒置連線數,沒在用就會關到剩這個數
	SqlMaxIdleSecond = mconfig.GetDuration("sql.maxIdleSecond")         //sql連線最大閒置秒數,超過閒置時間就會關閉連線
	Cq9GameHost = mconfig.GetString("thirdApi.cq9Games.host")           //cq9遊戲api host
	Cq9AuthToken = mconfig.GetString("thirdApi.cq9Games.authorization") //cq9 auth token
	Cq9GameHall = mconfig.GetString("thirdApi.cq9Games.gameHall")       //cq9遊戲廳
	Cq9TeamID = mconfig.GetString("thirdApi.cq9Games.teamId")           //cq9團隊代碼
	Mode = mconfig.GetString("application.mode")                        //程式環境,Integration:測試 Production:線上
	GameClientUrl = mconfig.GetString("api.gameClientUrl")              //遊戲客端網址,輸入gamecode/gametoken
	SystemDetailUrl = mconfig.GetString("api.systemDetailUrl")          //系統商遊戲細單客端網址,輸入gametoken
	PlayerDetailUrl = mconfig.GetString("api.playerDetailUrl")          //玩家端遊戲細單客端網址,輸入gametoken
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
