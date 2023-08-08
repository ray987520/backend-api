package uuid

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"fmt"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

/////////////////////////////
////    封裝的Uuid(sonyFlake)
/////////////////////////////

const (
	sonyFlakeBaseTime  = "2023-01-01 00:00:00.000"           //sonyflake default time
	flakeTimeFormat    = "2006-01-02 15:04:05.999"           //sonyflake default time format
	initFlakeTimeError = "init sonyflake base time error:%v" // init sonyflake time error message
	flakeInstanceError = "sonyflake instance error"          //sonyflake instance error message
)

var sonyFlake *sonyflake.Sonyflake

// 取機器ID
func getMachineID() (machineID uint16, err error) {
	//*TODO 暫時使用一個假的machineID,後續應有環境變數或其他方式提供機器ID
	machineID = 1688
	return machineID, nil
}

// 初始化,設置sonyFlake基礎值
func initSonyflake() {
	beginTime, err := time.Parse(flakeTimeFormat, sonyFlakeBaseTime)
	//parse time error,return
	if err != nil {
		err = fmt.Errorf(initFlakeTimeError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SonyflakeInit, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, err)
		return
	}
	st := sonyflake.Settings{
		StartTime: beginTime,
	}
	st.MachineID = getMachineID
	sonyFlake = sonyflake.NewSonyflake(st)
}

// 取得sonyflake初始後實體
func getInstance() *sonyflake.Sonyflake {
	//返回已初始實體
	if sonyFlake != nil {
		return sonyFlake
	}

	//否則初始化,再返回實體
	initSonyflake()

	return sonyFlake
}

// 產生sonyFlakeID
func Gen(traceCode string) (uuid string) {
	//無法取得實體,返回空字串
	if getInstance() == nil {
		err := fmt.Errorf(flakeInstanceError)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SonyflakeGenUuid, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		return ""
	}

	id, err := sonyFlake.NextID()
	//無法取得uuid,返回空字串
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SonyflakeGenUuid, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		return ""
	}

	//轉成16進位數字字串(比較短)
	uuid = strconv.FormatUint(id, 16)

	return uuid
}
