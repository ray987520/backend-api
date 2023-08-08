package str

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"strconv"
)

/////////////////////////////
////    封裝的字串轉換服務
/////////////////////////////

// 封裝strconv.Atoi
func Atoi(traceCode, input string) (data int, isOK bool) {
	data, err := strconv.Atoi(input)
	//String轉Int異常,返回int default:0跟失敗:false
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.StringConvertAtoi, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "input", input))
		return 0, false
	}

	return data, true
}

// 封裝strconv.Itoa
func Itoa(traceCode string, input int) (data string) {
	return strconv.Itoa(input)
}

// 封裝strconv.ParseFloat 64bit
func ParseFloat64(traceCode, input string) (data float64, isOK bool) {
	data, err := strconv.ParseFloat(input, 64)
	//String轉Float64異常,返回float64 default:0跟失敗:false
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.StringConvertParseFloat64, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "input", input))
		return 0, false
	}

	return data, true
}

// 封裝strconv.ParseBool
func ParseBool(traceCode, input string) (data bool, isOK bool) {
	data, err := strconv.ParseBool(input)
	//String轉Boolean異常,返回bool default:false跟失敗:false
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.StringConvertParseBool, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "input", input))
		return false, false
	}

	return data, true
}
