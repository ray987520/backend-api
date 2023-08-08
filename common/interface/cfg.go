package iface

import "time"

/////////////////////////////
////    Thirdparty Service Config Interface
/////////////////////////////

type IAppConfig interface {
	// get cq9遊戲api host
	GetCq9GameHost() string

	// get cq9 auth token
	GetCq9AuthToken() string

	// get cq9遊戲廳
	GetCq9GameHall() string

	// get cq9團隊代碼
	GetCq9TeamID() string

	// get sql master連線字串
	GetSqlConnectString() string

	// get sql最大連線數
	GetSqlMaxOpenConns() int

	// get sql最大閒置連線數
	GetSqlMaxIdleConns() int

	// get sql連線最大閒置秒數
	GetSqlMaxIdleSecond() time.Duration
}
