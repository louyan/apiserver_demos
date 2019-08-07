package middleware

import (
	"errors"
	"strings"
	"time"

	"yannotes.cn/apiserver_demos/demo06/handler"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		if token == "" {
			token = c.Request.Header.Get("Authorization")
			if s := strings.Split(token, " "); len(s) == 2 {
				token = s[1]
			}
		}

		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				if token, err := j.RefreshToken(token); err == nil {
					c.Header("Authorization", "Bear "+token)
					//c.JSON(http.StatusOK, gin.H{"error": 0, "message": "refresh token", "token": token})
					handler.SendResponse(c, nil, token)
					return
				}
			}
			//c.JSON(http.StatusUnauthorized, gin.H{"error": 1, "message": err.Error()})
			handler.SendResponse(c, TokenUnauthorized, nil)
			return
		}

		c.Set("claims", claims)
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired            = errors.New("Token is expired")
	TokenNotValidYet  error = errors.New("Token not active yet")
	TokenMalformed    error = errors.New("That's not even a token")
	TokenInvalid      error = errors.New("Couldn't handle this token:")
	TokenUnauthorized       = errors.New("Unauthorized.")
)

type CustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username`
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

func GetSignKey() string {
	return viper.GetString("jwt_secret")
}

func SetSignKey(key string) string {
	viper.Set("jwt_secret", key)
	return GetSignKey()
}

//GenToken 生成token
func (j *JWT) GenToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 { // Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 { // TokenNotValidYet
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

//RefreshToken 刷新Token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", nil
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.GenToken(*claims)
	}

	return "", TokenInvalid
}
