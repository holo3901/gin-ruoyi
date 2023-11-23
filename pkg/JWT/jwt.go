package JWT

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/spf13/viper"
)

var mySecret = []byte("miku") //定义加密
// MyClaims 自定义声明结构体并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenToken(cardID int64, username string) (aToken string, err error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		cardID,
		username, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * viper.GetDuration("auth.jwt_expire"))},
			Issuer:    "ruoyi", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象	// 使用指定的secret签名并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	//refresh token不需要存任何自定义数据
	/*rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Second * 30)},
		Issuer:    "bluebell",
	}).SignedString(mySecret)*/
	return
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var token *jwt.Token
	//way1将解析结果保存到claims变量中,若token字符串合法但过期claims也会有数据，err会提示token过期
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { //校验token
		return mc, nil
	}
	/*//way2从ParseWithClaims返回的Token结构体中取出Claims结构体
	token, err := jwt.ParseWithClaims(tokenString, *MyClaims, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	*/
	return nil, errors.New("invalid token")
}

/*
//RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	//refresh token无效直接返回
	if _, err = jwt.Parse(rToken, func(aToken *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	}); err != nil {
		return
	}
	//从旧access token中解析出claims数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, func(aToken *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	v, _ := err.(*jwt.ValidationError)
	//当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Username)
	}
	return
}
*/
