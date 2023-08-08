package errorcode

/////////////////////////////
////    API錯誤代碼
/////////////////////////////

type ErrorCode string

// api錯誤代碼
const (
	Success      ErrorCode = "0"
	BackendError ErrorCode = "999"
)
