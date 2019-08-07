package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

var (
	ErrMissingHeader = errors.New("The length of the `Authorization` header is zero.")
)

type Context struct {
	ID       uint64
	Username string
}

//Signature 签名,secret加密
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

//用密钥解析token，验证有效则返回响应数据
func Parse(tokenString, secret string) (*Context, error) {
	ctx := &Context{}

	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil {
		return ctx, nil

		// Read the token if it's valid
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.ID = uint64(claims["id"].(float64))
		ctx.Username = claims["username"].(string)
		return ctx, nil
	} else {
		return ctx, err
	}
}

//ParseRequest 解析请求，获取token
func ParseRequest(c *gin.Context) (*Context, error) {
	header := c.Request.Header.Get("Authorization")
	fmt.Println("=======", header)
	//Load the jwt secret from config
	secret := viper.GetString("jwt_secret")

	if len(header) == 0 {
		return &Context{}, ErrMissingHeader
	}

	var t string

	// Parse the header to get the token part
	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, secret)
}

//签发 token
func Sign(ctx *gin.Context, c Context, secret string) (tokenString string, err error) {
	if secret == "" {
		secret = viper.GetString("jwt_secret")
	}

	// The token content
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       c.ID,
		"username": c.Username,
		"nbf":      time.Now().Unix(), //JWT Token 生效时间
		"iat":      time.Now().Unix(), //JWT Token 签发时间
	})

	//使用指定的密钥对令牌签名
	tokenString, err = token.SignedString([]byte(secret))

	return
}
