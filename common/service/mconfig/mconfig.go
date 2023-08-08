package mconfig

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

/////////////////////////////
////    封裝的Config服務(viper)
/////////////////////////////

const (
	viperReadFileError   = "viper read config file error:%v"
	viperReadConfigError = "viper read config error ,configPath:%s ,data:%v"
	configChangeMessage  = "config file changed ,data:%s"
	configPath           = "./cfg/"
)

// 初始化viper,採用初始化時傳入config檔名供不同包API使用
func InitConfigManager(configFileName string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configFileName)
	err := viper.ReadInConfig()
	//讀取設定檔失敗,不繼續執行代碼
	if err != nil {
		err = fmt.Errorf(viperReadFileError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.MConfigInitConfigManager, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, err)
		panic(err)
	}
	//動態監看設定檔更新
	viper.WatchConfig()
	//設定檔更新時列印更新欄位
	viper.OnConfigChange(func(e fsnotify.Event) {
		msg := fmt.Sprintf(configChangeMessage, e.Name)
		zaplog.Infow(innertrace.InfoNode, innertrace.FunctionNode, thirdparty.MConfigInitConfigManager, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, msg)
	})
}

// 取設定值,string
func GetString(configPath string) string {
	return cast.ToString(Get(configPath))
}

// 取設定值,int
func GetInt(configPath string) int {
	return cast.ToInt(Get(configPath))
}

// 取設定值,int64
func GetInt64(configPath string) int64 {
	return cast.ToInt64(Get(configPath))
}

// 取設定值,時間區間
func GetDuration(configPath string) time.Duration {
	return cast.ToDuration(Get(configPath))
}

// 取設定值,[]string
func GetStringSlice(configPath string) []string {
	return cast.ToStringSlice(Get(configPath))
}

// 取設定值,interface{}
func Get(configPath string) any {
	data := viper.Get(configPath)
	//如果找不到設定值,不希望代碼繼續執行
	if data == nil {
		err := fmt.Errorf(viperReadConfigError, configPath, data)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.MConfigGet, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, err)
		panic(err)
	}
	return data
}
