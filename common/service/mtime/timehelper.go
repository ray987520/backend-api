package mtime

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"encoding/json"
	"strings"
	"time"
)

/////////////////////////////
////    封裝的時間共用服務
/////////////////////////////

const (
	ApiTimeFormat           = "2006-01-02T15:04:05.999-07:00"  //文件定義API時間格式
	SysTimeFormat           = "2006-01-02 15:04:05.999"        //system的時間格式
	DbTimeFormat            = "2006-01-02 15:04:05.999 -07:00" //sql db吃的時間格式
	LogTimeFormat           = "2006-01-02 15:04:05.999999"     //log的時間格式,精準度較高
	BackendUtcTimeFormat    = "2006-01-02 15:04:05"            //後台Utc時間格式
	BackendReportTimeFormat = "2006-01-02T15:04:05.999Z"       //後台爬蟲報表時間格式
	DateTimeOffsetMinValue  = "0001-01-01 00:00:00.000 +00:00" //DateTimeOffset.MinValue(1/1/0001 12：00：00 AM +00：00)
)

var (
	defaultTimeZone int       //default時區
	Default         time.Time //default時間(1/1/0001 12：00：00 AM +00：00)
)

// 初始化封裝的時間服務
func InitTimeService(timezone int) {
	defaultTimeZone = timezone
	Default, _ = time.Parse(DbTimeFormat, DateTimeOffsetMinValue)
}

// UTC時間
func UtcNow() time.Time {
	return time.Now().UTC()
}

// 本地時間,依設定值api.timezone計算
func LocalNow() time.Time {

	zone := GetTimeZone(defaultTimeZone)
	return UtcNow().In(zone)
}

// 依format轉出時間字串
func TimeString(t time.Time, format string) string {
	return t.Format(format)
}

// 依format轉出時間字串,較短的話補0,ex:2023-07-01 00:00:00.12=>2023-07-01 00:00:00.120
func TimeStringAndFillZero(t time.Time, format string) string {
	data := t.Format(format)
	//如果比較短才處理
	if len(data) < len(format) {
		t = t.Add(1 * time.Millisecond)
		data = t.Format(format)
		data = strings.TrimRight(data, "1") + "0"
	}
	return data
}

// 按format parse UTC0時間
func ParseTime(traceCode, format, timeString string) (t time.Time, isOK bool) {
	t, err := time.Parse(format, timeString)
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.ParseTime, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "format", format, "timeString", timeString))
		return t, false
	}
	return t, true
}

// 按format parse UTC+default timeZone時間
func ParseToLocal(traceCode, format, timeString string) (t time.Time, isOK bool) {
	return ParseToTimeZone(traceCode, format, timeString, defaultTimeZone)
}

// 按format parse UTC+target timeZone時間
func ParseToTimeZone(traceCode, format, timeString string, targetTimeZone int) (t time.Time, isOK bool) {
	zone := GetTimeZone(targetTimeZone)
	t, isOK = ParseTime(traceCode, format, timeString)
	if !isOK {
		return t, false
	}
	return t.In(zone), true
}

// 自訂轉ApiTime成SystemTime時間類型,ex:2023-07-01T00:00:00.120-04:00=>2023-07-01 00:00:00.120
type ApiTime time.Time

// 自訂轉ApiTime成SystemTime時間類型實作UnmarshalJSON
func (d *ApiTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(ApiTimeFormat, s)
	if err != nil {
		return err
	}
	*d = ApiTime(t)
	return nil
}

// 自訂轉ApiTime成SystemTime時間類型實作MarshalJSON
func (d ApiTime) MarshalJSON() ([]byte, error) {
	s := TimeStringAndFillZero(time.Time(d), SysTimeFormat)
	return json.Marshal(s)
}

// 取特定timezone,ex:+08:00=>GetTimeZone(8)
func GetTimeZone(areaPower int) *time.Location {
	zone := time.FixedZone("", areaPower*60*60)
	return zone
}
