package innertrace

/////////////////////////////
////    Log用的共用節點名稱
/////////////////////////////

//log節點名稱
const (
	TraceNode       = "x-game-tracecode"   //用於zaplog,traceCode節點名稱
	FunctionNode    = "function"           //用於zaplog,function節點名稱
	ErrorInfoNode   = "error"              //用於zaplog,error節點名稱
	ErrorCodeNode   = "errorcode"          //用於zaplog,errorcode節點名稱
	DataNode        = "content"            //用於zaplog,data節點名稱
	InfoNode        = "processInfo"        //用於zaplog,info節點名稱
	RequestTimeNode = "x-game-requesttime" //用於zaplog,requesttime節點名稱
	TotalTimeNode   = "totalTime"          //用於zaplog,totalTime節點名稱
)

//log summary
const (
	DBSqlError            = "sql error"                          //用於zaplog,log summary
	ExternalServiceError  = "third party service error"          //用於zaplog,log summary
	ParseFormError        = "parse formdata error"               //用於zaplog,log summary
	ValidRequestError     = "bad http request"                   //用於zaplog,log summary
	ConfigError           = "config error"                       //用於zaplog,log summary
	PanicError            = "panic error"                        //用於zaplog,log summary
	MiddlewareError       = "middleware error"                   //用於zaplog,log summary
	LogOriginRequest      = "log original request"               //用於zaplog,log summary
	LogRequestModel       = "log parsed request model"           //用於zaplog,log summary
	LogThirdApiRequest    = "log call thirdparty api request"    //用於zaplog,log summary
	LogThirdApiParameters = "log call thirdparty api parameters" //用於zaplog,log summary
	LogThirdApiResponse   = "log call thirdparty api response"   //用於zaplog,log summary
	LogValidate           = "log self define validate"           //用於zaplog,log summary
	ParseCq9TimeError     = "parse cq9 status.time error"        //用於zaplog,log summary
)
