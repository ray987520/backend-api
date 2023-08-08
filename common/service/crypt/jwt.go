package crypt

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	"GamePoolApi/common/service/mtime"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/uuid"
	"GamePoolApi/common/service/zaplog"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

/////////////////////////////
////    JWT token加解密服務
/////////////////////////////

// JWT token要加密的資料
type AccessTokenClaims struct {
	Id        string `json:"jti"`       // jwt token id,uuid
	Type      string `json:"typ"`       // jwt token類型,JSON Web Token
	Sub       string `json:"sub"`       // 放帳號,admin
	TimeZone  string `json:"TimeZone"`  // 時區
	LoginID   string `json:"LoginID"`   // 表示 MemberID ，使用者識別碼
	Authority string `json:"Authority"` // 表示使用者權限
	NotBefore int64  `json:"nbf"`       // Token 在什麼時間之前生效
	ExpiresAt int64  `json:"exp"`       //Token 的逾期時間
	IssuedAt  int64  `json:"iat"`       // Token 的建立時間
	Issuer    string `json:"iss"`       // 發送 Token 的發行者
}

const (
	issuer             = "GenesisPerfectGame"       //jwt issuer
	authority          = "[1,2,4,5,7,8,6,5]"        //使用者權限
	signStringError    = "gen token error:%v"       //jwt sign error
	parseJwtError      = "parse jwt token error:%v" //parse jwt error
	getJwtDataError    = "get data of jwt error:%v" //get jwt data error
	expireMinutes      = 480                        //jwt token過期時間,480 minutes
	jwtType            = "JSON Web Token"
	expireErrorMessage = "token expired"
)

var (
	jwtSecret []byte // jwt secret key
)

// 初始化帶入JWT密鑰
func InitJwt(jwtKey string) {
	jwtSecret = []byte(jwtKey)
}

// 產生JWT TOKEN
func JwtCreateAccessToken(traceCode string, memberId string, account, timeZone string) (tokenString string) {
	zaplog.Infow(innertrace.InfoNode, innertrace.FunctionNode, thirdparty.JwtCreateAccessToken, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("account", account, "timeZone", timeZone, "memberId", memberId, "jwtSecret", jwtSecret))

	//設定jwt要加密的內容
	now := mtime.UtcNow()
	claims := new(AccessTokenClaims)
	claims.Id = uuid.Gen(traceCode)
	claims.Type = jwtType
	claims.Sub = account
	claims.LoginID = memberId
	claims.TimeZone = timeZone
	claims.Authority = authority
	claims.NotBefore = now.Unix()
	claims.ExpiresAt = now.Add(expireMinutes * time.Minute).Unix()
	claims.IssuedAt = now.Unix()
	claims.Issuer = issuer

	//簽章加密選用SHA512
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	// jwtToken:=&jwt.Token{
	// 		Header: map[string]interface{}{
	// 			"typ": "JWT",
	// 			"alg": jwt.SigningMethodHS512.Alg(),
	// 		},
	// 		Claims: claims,
	// 		Method: jwt.SigningMethodHS512,
	// 	}
	tokenString, err := jwtToken.SignedString(jwtSecret)
	//jwt sign error返回空值
	if err != nil {
		err = fmt.Errorf(signStringError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.JwtCreateAccessToken, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, err, "tokenString", tokenString))
		return ""
	}

	return tokenString
}

// 驗證JWT TOKEN
func JwtValidAccessToken(traceCode string, tokenString string) *AccessTokenClaims {
	zaplog.Infow(innertrace.InfoNode, innertrace.FunctionNode, thirdparty.JwtValidAccessToken, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage("tokenString", tokenString, "jwtSecret", jwtSecret))

	//若傳進來還是Bearer Auth的樣式,處理掉Bearer頭
	if strings.HasPrefix(tokenString, "Bearer") {
		tokenString = strings.Split(tokenString, " ")[1]
	}

	//解析jwt token
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	//parse jwt error返回空值
	if err != nil {
		err = fmt.Errorf(parseJwtError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.JwtValidAccessToken, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
		return nil
	}

	// 從raw token中取回資訊,成功就返回資訊
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims
	}

	//無法成功取回對應格式資料就是jwt字串異常
	err = fmt.Errorf(getJwtDataError, err)
	zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.JwtValidAccessToken, innertrace.TraceNode, traceCode, innertrace.DataNode, err)
	return nil
}

// AccessTokenClaims驗證,若完全自訂沒使用套件的standard claim就必須實作
func (claim *AccessTokenClaims) Valid() error {
	now := mtime.UtcNow().Unix()
	//若token.ExpiresAt超時返回錯誤
	if claim.ExpiresAt < now {
		return fmt.Errorf(expireErrorMessage)
	}
	return nil
}
