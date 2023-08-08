package crypt

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"encoding/base64"
)

/////////////////////////////
////    base64編碼服務
/////////////////////////////

// base64解碼
func Base64Decode(traceCode, data string) ([]byte, bool) {
	zaplog.Infow(innertrace.InfoNode, innertrace.FunctionNode, thirdparty.Base64Decode, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("data", data))

	decode, err := base64.StdEncoding.DecodeString(data)
	//解碼錯誤返回null
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.Base64Decode, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err))
		return nil, false
	}

	return decode, true
}
