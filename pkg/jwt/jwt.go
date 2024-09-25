package jwt

import (
	"errors"
	"project_bluebell/settings"
	"time"

	"github.com/golang-jwt/jwt"
)

/*
   jwt :   astring.bstring.cstring  头部.负载.签名
   头部: json对象，存放jwt元数据，编码后转成字符串放入astring
   负载: json对象，存放实际需要传递的数据，，同样编码后转成字符串放入bstring
   签名: 签名是对前两部分的签名，先指定一个密钥，根据指定的签名算法产生一个字符串，放入cstring
   将 头部，负载，签名三部分用.拼接成字符串，就是 jwt ， 返回给用户

   优点: 默认不加密，但是可以加密。可以降低服务器查询数据库的次数，
   缺点: 在有效时间内无法废弃终止某个token，一旦泄露，任何人都可以获得这个token 的权限。所以通常将有效时间设置得很短
*/

var ErrorInvalidToken = errors.New("invalid token")

var MySecret = []byte("encryption")

// 自定义 jwt 保存的数据字段
type MyClaims struct {
	User_ID            int64  `json:"user_id"`
	Username           string `json:"username"`
	jwt.StandardClaims        // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(user_id int64, username string) (string, error) {

	// 创建一个我们自己的声明
	claims := MyClaims{
		user_id,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(settings.Conf.AuthConfig.Jwt_time) * time.Hour).Unix(), // 过期时间
			Issuer:    "Bluebell",                                                                          // 签发人
		},
	}

	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
// 把用户登录是传给服务器的 tokenString 解析成我们定义的结构体类型，里面可以拿到用户 ID 和 Name
func ParseToken(tokenString string) (*MyClaims, error) {

	// 解析token
	var mc = new(MyClaims) // 这里需要手动初始化一个内存
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})

	if err != nil {
		return nil, err
	}
	// 校验 token
	if token.Valid {
		return mc, nil
	}
	return nil, ErrorInvalidToken
}

// Refresh Token 刷新 access token
// *这里的返回值是可以根据实际功能来决定的
//
//	-1. 同时返回 accessToken 和 refreshToken , 此时，只要在 refresh token 的有效期内进行再次登录就不会重新登录
//	-2. 只返回刷新过的 accessToken , 此时，只要过了 refresh token 的有效期就必须重新登陆
func RefreshToken(aToken string, rToken string) (newAToken string, err error) {

	// refresh token 过期直接返回
	// jet.Parse(token , keyFunc) 里面的函数就是下面传入的函数，返回密钥和错误
	_, err = jwt.Parse(rToken, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})

	if err != nil {
		return "", err
	}

	// 从原来的 access token 中解析出原来的数据
	var claims MyClaims
	token, err := jwt.ParseWithClaims(aToken, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})

	// 如果 access token 是过期错误，并且 refresh token 没有过期，就重新生成 token 并返回
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors == jwt.ValidationErrorExpired {
				return GenToken(claims.User_ID, claims.Username)
			}
		}
		return "", err
	}

	// aToken 没有过期，可以直接返回，也可以重新生成
	if token.Valid {
		return aToken, nil
	}

	return "", err
}
