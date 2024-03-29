package jwt

import (
	"duoduo_fun/pkg/errno"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var mySecret = []byte("Jiyeon_Hyomin_Jiyeon_Hyomin_hhhh")

// MyClaims CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段s
// 假设我们这里需要额外记录一个userid字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	// 可根据需要自行添加字段
	UserID               int    `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(userid int, username string) (aToken string, err error) {
	// 实例化一个我们带创建的加密声明
	aclaims := MyClaims{
		// 自定义字段
		userid,
		username,
		jwt.RegisteredClaims{
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			// 签发人
			Issuer: "朴智妍",
		},
	}

	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, aclaims).SignedString(mySecret)
	return
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = errno.NewErrNo("Token解析出错" + jwt.ErrTokenExpired.Error())
		}
		return nil, err
	}
	// 校验token
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
