package middleware

/////////////////////////////
////    Middleware Function Name
/////////////////////////////

type MiddlewareId string

// middleware function,用於zaplog分類
const (
	AuthMiddleware        MiddlewareId = "AuthMiddleware"
	AcceptMiddleware      MiddlewareId = "AcceptMiddleware"
	SelfHeaderMiddleware  MiddlewareId = "SelfHeaderMiddleware"
	ErrorHandleMiddleware MiddlewareId = "ErrorHandleMiddleware"
	IPWhiteListMiddleware MiddlewareId = "IPWhiteListMiddleware"
	LogOriginRequest      MiddlewareId = "logOriginRequest"
	TotalTimeMiddleware   MiddlewareId = "TotalTimeMiddleware"
	BasicAuth             MiddlewareId = "basicAuth"
)
