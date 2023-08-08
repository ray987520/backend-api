package zaplog

import (
	"GamePoolApi/common/enum/reqireheader"
	"GamePoolApi/common/service/notify"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/////////////////////////////
////    封裝的Log服務(zap+lumberjack)
/////////////////////////////

type Level zapcore.Level

// log層級
const (
	DebugLevel Level = Level(zap.DebugLevel)
	InfoLevel  Level = Level(zap.InfoLevel)
	ErrorLevel Level = Level(zap.ErrorLevel)
)

// var會比init早執行,但zaplog這裡不是用mconfig,所以先給預設值
var (
	maxlogsize      int
	maxbackup       int
	maxage          int
	svcname         string
	logFilePath     string
	defaultLogLevel string
)

var (
	logger      *zap.SugaredLogger //log實例
	atomicLevel = zap.NewAtomicLevel()

	levelMap = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"error": zapcore.ErrorLevel,
	}
	zviper *viper.Viper
)

const (
	viperReadFileError   = "viper read config file error:%v"                 //read config file error message
	viperReadConfigError = "viper read config error ,configPath:%s ,data:%v" //read setting error message
	configChangeMessage  = "config file changed ,data:%s"                    //monitor config setting change message
	configPath           = "./cfg/"
)

// 初始化zaplog
func InitZaplog(configFileName string) {
	//因為zaplog為最底層,跟mconfig又不同包,避免循環參照只能拉到同一層或是獨立viper
	//viper讀取config
	zviper = viper.New()
	zviper.AddConfigPath(configPath)
	zviper.SetConfigName(configFileName)
	err := zviper.ReadInConfig()
	//read config file error,panic error
	if err != nil {
		panic(fmt.Sprintf(viperReadFileError, err))
	}

	//初始化zaplog實體,目前選用sugar log,跟原生log封裝比較相近,但效能較低
	maxlogsize = zviper.GetInt("log.maxlogsize")              //每50MB切割log檔
	maxbackup = zviper.GetInt("log.maxbackup")                //達300個切割開始取代舊檔
	maxage = zviper.GetInt("log.maxage")                      //log保存最大時間(days)
	svcname = zviper.GetString("log.svcname")                 //log檔名
	logFilePath = zviper.GetString("log.logFilePath")         //log檔案路徑
	defaultLogLevel = zviper.GetString("log.defaultLogLevel") //預設log level
	filePath := getFilePath()
	level := getLoggerLevel(defaultLogLevel)
	log := NewLogger(filePath, level, maxlogsize, maxbackup, maxage, true, svcname)
	logger = log.Sugar()
	logger.Sync()

	//初始化line警示服務
	notify.InitLineNotify(fmt.Sprintf(reqireheader.JwtTokenFormat, zviper.GetString("thirdApi.LineNotify.Token"))) //line警示access token,格式為"Bearer XXXXXXXXXXXXXXX"
}

// 取logger層級,預設info
func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// 用於動態調整logger層級
// *TODO 提供一個API便於程式執行時切換logger層級
func SetLoggerLevel(lvl string) {
	filePath := getFilePath()
	level := getLoggerLevel(lvl)
	log := NewLogger(filePath, level, maxlogsize, maxbackup, maxage, true, svcname)
	logger = log.Sugar()
	logger.Sync()
}

// 建立zap的實例,添加caller顯示調用log的代碼位置
func NewLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, serviceName string) *zap.Logger {
	core := newCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddCallerSkip(2))
}

// 設置zapcore,添加輸出到os.Stdout跟lumberjack log切割
func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,                                       //行結尾
		EncodeLevel:    zapcore.CapitalLevelEncoder,                                     // 全大寫
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:06:999999Z07:00"), //zapcore.RFC3339NanoTimeEncoder,     // RFC3339 UTC 时间格式,2006-01-02T15:04:05.999999999Z07:00
		EncodeDuration: zapcore.SecondsDurationEncoder,                                  //顯示浮點樹的秒
		EncodeCaller:   zapcore.ShortCallerEncoder,                                      // 短調用路徑
		// EncodeCaller:   zapcore.FullCallerEncoder,    // 長調用路徑
		EncodeName: zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		//設定同步寫入log的標的
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout), // 打印到控制台
			zapcore.AddSync(&hook)),    // 打印到文件
		zap.NewAtomicLevelAt(level), // 日志级别
	)
}

// 取API執行的當前目錄
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Error(err)
	}
	return dir
}

// 取設定的log檔案路徑
func getFilePath() string {
	logfile := fmt.Sprintf(logFilePath, svcname)
	return logfile
}

// 輸出debug log,單一參數轉json
func Debug(arg interface{}) {
	data, err := json.Marshal(arg)
	if err != nil {
		logger.Debug(arg)
		return
	}
	logger.Debug(string(data))
}

// 輸出debug log,類似fmt.Sprintf
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// 輸出debug log,zap的keyvalue風格,配合一些log portal查詢較明確不須全文本查詢
func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

// 輸出info log,單一參數轉json
func Info(arg interface{}) {
	data, err := json.Marshal(arg)
	if err != nil {
		logger.Info(arg)
		return
	}
	logger.Info(string(data))
}

// 輸出info log,類似fmt.Sprintf
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// 輸出info log,zap的keyvalue風格,配合一些log portal查詢較明確不須全文本查詢
func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

// 輸出error log,單一參數轉json
func Error(arg interface{}) {
	data, err := json.Marshal(arg)
	//無法parse json的話直接log值
	if err != nil {
		// //send line alert
		// notify.Send(fmt.Sprintf("%v", arg))

		logger.Error(arg)
		return
	}

	// //send line alert
	// notify.Send(string(data))

	logger.Error(string(data))
}

// 輸出error log,類似fmt.Sprintf
func Errorf(template string, args ...interface{}) {
	// //send line alert
	// notify.Send(fmt.Sprintf(template, args...))

	logger.Errorf(template, args...)
}

// 輸出error log,zap的keyvalue風格,配合一些log portal查詢較明確不須全文本查詢
func Errorw(msg string, keysAndValues ...interface{}) {
	// //send line alert
	// notify.Send(mergeMessage("message", msg, keysAndValues))

	logger.Errorw(msg, keysAndValues...)
}

// 合併資料成log內容字串
func mergeMessage(keyValues ...interface{}) string {
	//input資料長度不正確返回空值
	if len(keyValues) == 0 || len(keyValues)%2 == 1 {
		return ""
	}

	//成對處理資料
	var messages []string
	for i := 0; i < len(keyValues); i += 2 {
		messages = append(messages, fmt.Sprintf("%v:%v", keyValues[i], keyValues[i+1]))
	}

	return strings.Join(messages, ", ")
}
