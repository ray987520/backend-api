package validator

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"regexp"

	"github.com/go-playground/validator/v10"
)

/////////////////////////////
////    封裝的驗證器服務(validate)
/////////////////////////////

var validate *validator.Validate

// 初始化驗證器並註冊自訂驗證器
func init() {
	//*TODO 若有需要其他語言錯誤訊息也需要在此添加翻譯元件
	validate = validator.New()
	validate.RegisterValidation("acct", ValidateAccount)          //自訂帳號驗證器
	validate.RegisterValidation("amt", ValidateAmount)            //自訂金額驗證器
	validate.RegisterValidation("mtcode", ValidateMtCode)         //自訂MtCode驗證器
	validate.RegisterValidation("farray", ValidateFloatArray)     //自訂FloatArray驗證器
	validate.RegisterValidation("utc", ValidateUtc)               //自訂UtcTime驗證器
	validate.RegisterValidation("mbacct", ValidateBackendAccount) // 自訂ManageBackground帳號驗證器
	validate.RegisterValidation("tz", ValidateTimeZone)           // 自訂TimeZone驗證器
}

// 驗證結構,按struct的validate tag
func IsValidStruct(traceCode string, data interface{}) bool {
	err := validate.Struct(data)
	//validate套件驗證到錯誤
	if err != nil {
		zaplog.Errorw(innertrace.ValidRequestError, innertrace.FunctionNode, thirdparty.IsValidStruct, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		return false
	}

	return err == nil
}

// 自訂帳號驗證器,CQ9玩家帳號字串長度限制36個字元
func ValidateAccount(f1 validator.FieldLevel) bool {
	match, err := regexp.MatchString("^[a-zA-Z0-9_-]{1,36}$", f1.Field().String())
	if err != nil {
		zaplog.Errorw(innertrace.ValidRequestError, innertrace.FunctionNode, thirdparty.ValidateAccount, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, tracer.MergeMessage("fieldName", f1.FieldName(), "fieldValue", f1.Field().String(), innertrace.ErrorInfoNode, err))
		return false
	}
	return match
}

// 自訂金額驗證器,最大長度為12位數及小數點後4位
func ValidateAmount(f1 validator.FieldLevel) bool {
	match, err := regexp.MatchString("^[0-9]{1,7}(.[0-9]{1,4})?$", f1.Field().String()) //1234567.1234
	if err != nil {
		zaplog.Errorw(innertrace.ValidRequestError, innertrace.FunctionNode, thirdparty.ValidateAmount, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, tracer.MergeMessage("fieldName", f1.FieldName(), "fieldValue", f1.Field().String(), innertrace.ErrorInfoNode, err))
		return false
	}
	return match
}

// 自訂MtCode驗證器,MTcode: {env}-{action}-{roundid},RoundID:團隊代碼+任意英文字母或數字,RoundID總長度不可以超過30碼
func ValidateMtCode(f1 validator.FieldLevel) bool {
	match, err := regexp.MatchString("^(rel|pro)-(bet|win|rollout|rollin)-[a-zA-Z0-9]{1,30}$", f1.Field().String()) //rel-bet-CC1213ab
	if err != nil {
		zaplog.Errorw(innertrace.ValidRequestError, innertrace.FunctionNode, thirdparty.ValidateMtCode, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, tracer.MergeMessage("fieldName", f1.FieldName(), "fieldValue", f1.Field().String(), innertrace.ErrorInfoNode, err))
		return false
	}
	return match
}

// 自訂FloatArray驗證器
func ValidateFloatArray(f1 validator.FieldLevel) bool {
	match, err := regexp.MatchString(`^\[([0-9]{1,7}(.[0-9]{1,4})?){1,}(,[0-9]{1,7}(.[0-9]{1,4})?)?\]$`, f1.Field().String()) //[1.1,2.2]
	if err != nil {
		zaplog.Errorw(innertrace.ValidRequestError, innertrace.FunctionNode, thirdparty.ValidateFloatArray, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, tracer.MergeMessage("fieldName", f1.FieldName(), "fieldValue", f1.Field().String(), innertrace.ErrorInfoNode, err))
		return false
	}
	return match
}

// 自訂UtcTime驗證器
func ValidateUtc(f1 validator.FieldLevel) bool {
	match, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`, f1.Field().String()) //2016-01-19T15:21:32.591Z
	if err != nil {
		zaplog.Errorw(innertrace.ValidRequestError, innertrace.FunctionNode, thirdparty.ValidateUtc, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, tracer.MergeMessage("fieldName", f1.FieldName(), "fieldValue", f1.Field().String(), innertrace.ErrorInfoNode, err))
		return false
	}
	return match
}

// 自訂ManageBackground帳號驗證器,ManageBackground帳號字串長度限制2~30個字元
func ValidateBackendAccount(f1 validator.FieldLevel) bool {
	match, err := regexp.MatchString("^[a-zA-Z0-9_-]{2,30}$", f1.Field().String())
	if err != nil {
		zaplog.Errorw(innertrace.ValidRequestError, innertrace.FunctionNode, thirdparty.ValidateBackendAccount, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, tracer.MergeMessage("fieldName", f1.FieldName(), "fieldValue", f1.Field().String(), innertrace.ErrorInfoNode, err))
		return false
	}
	return match
}

// 自訂TimeZone驗證器
func ValidateTimeZone(f1 validator.FieldLevel) bool {
	match, err := regexp.MatchString(`^[+-]\d{2}[:]\d{2}$`, f1.Field().String()) //+08:00
	if err != nil {
		zaplog.Errorw(innertrace.ValidRequestError, innertrace.FunctionNode, thirdparty.ValidateTimeZone, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, tracer.MergeMessage("fieldName", f1.FieldName(), "fieldValue", f1.Field().String(), innertrace.ErrorInfoNode, err))
		return false
	}
	return match
}
