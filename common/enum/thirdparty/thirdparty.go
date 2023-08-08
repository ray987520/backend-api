package thirdparty

/////////////////////////////
////    第三方或封裝服務 Function Name
/////////////////////////////

type FunctionId string

// thirdparty function,用於zaplog
const (
	ElasticInit                 FunctionId = "initElastic"
	ElasticCreateIndex          FunctionId = "createElasticIndex"
	ElasticGenUuid              FunctionId = "genElasticUuid"
	ParseTime                   FunctionId = "ParseTime"
	ParseToLocal                FunctionId = "ParseToLocal"
	Aes128Encrypt               FunctionId = "Aes128Encrypt"
	Aes128Decrypt               FunctionId = "Aes128Decrypt"
	JwtCreateAccessToken        FunctionId = "JwtCreateAccessToken"
	JwtValidAccessToken         FunctionId = "JwtValidAccessToken"
	SonyflakeInit               FunctionId = "InitSonyflake"
	SonyflakeGenUuid            FunctionId = "Gen"
	JsonMarshal                 FunctionId = "JsonMarshal"
	JsonUnMarshal               FunctionId = "JsonUnMarshal"
	StringConvertAtoi           FunctionId = "Atoi"
	StringConvertParseFloat64   FunctionId = "ParseFloat64"
	StringConvertParseBool      FunctionId = "ParseBool"
	GormInit                    FunctionId = "gormInit"
	SqlSelect                   FunctionId = "Select"
	SqlUpdate                   FunctionId = "Update"
	SqlDelete                   FunctionId = "Delete"
	SqlCreate                   FunctionId = "Create"
	SqlTransaction              FunctionId = "Transaction"
	MConfigInitConfigManager    FunctionId = "InitConfigManager"
	MConfigGet                  FunctionId = "Get"
	MConfigGetString            FunctionId = "GetString"
	MConfigGetInt               FunctionId = "GetInt"
	MConfigGetInt64             FunctionId = "GetInt64"
	MConfigGetDuration          FunctionId = "GetDuration"
	MConfigGetStringSlice       FunctionId = "GetStringSlice"
	GetErrorHttpResponse        FunctionId = "getErrorHttpResponse"
	GetBackendErrorHttpResponse FunctionId = "getBackendErrorHttpResponse"
	HttpRequest2Curl            FunctionId = "HttpRequest2Curl"
	WriteHttpResponse           FunctionId = "WriteHttpResponse"
	CallThirdApi                FunctionId = "CallThirdApi"
	CallThirdApiByBody          FunctionId = "CallThirdApiByBody"
	ReadHttpRequestBody         FunctionId = "ReadHttpRequestBody"
	ToQueryString               FunctionId = "ToQueryString"
	IsValidStruct               FunctionId = "IsValidStruct"
	ValidateAccount             FunctionId = "ValidateAccount"
	ValidateAmount              FunctionId = "ValidateAmount"
	ValidateMtCode              FunctionId = "ValidateMtCode"
	ValidateFloatArray          FunctionId = "ValidateFloatArray"
	ValidateUtc                 FunctionId = "ValidateUtc"
	ValidateBackendAccount      FunctionId = "ValidateBackendAccount"
	ValidateTimeZone            FunctionId = "ValidateTimeZone"
	Base64Decode                FunctionId = "Base64Decode"
)
