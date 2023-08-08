package reqireheader

/////////////////////////////
////    API必須的request/response header
/////////////////////////////

//api必須的request header
const (
	Authorization    = "Authorization"     //認證header
	ContentType      = "content-type"      //文本格式header
	CfConnectingIp   = "CF-Connecting-IP"  //取remote ip最優先header
	XForwardedFor    = "X-Forwarded-For"   //取remote ip次優先header
	TransferEncoding = "Transfer-Encoding" //原API response固定header

	FormUrlEncode = "application/x-www-form-urlencoded" //http header value ofcontent-type
	Json          = "application/json"                  //http header value ofcontent-type
	Chunked       = "chunked"                           //原API response固定header Transfer-Encoding值

	JwtTokenFormat = "Bearer %s"
)
