package serializer

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"encoding/json"
)

/////////////////////////////
////    封裝的Json序列化反序列化
/////////////////////////////

// 封裝Json序列化
func JsonMarshal(traceCode string, v any) (data []byte) {
	data, err := json.Marshal(v)
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.JsonMarshal, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "v", v))
		return nil
	}
	return data
}

// 封裝Json反序列化,v請傳址
func JsonUnMarshal(traceCode string, data []byte, v any) (isOK bool) {
	err := json.Unmarshal(data, v)
	if err != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.JsonUnMarshal, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "data", string(data)))
		return false
	}
	return true
}
