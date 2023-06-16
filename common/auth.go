package common

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TaoClaims struct {
	Username string `json:"username"`
	OpenId   string `json:"openid"`
	Randomid string `json:"randomid"`
	jwt.RegisteredClaims
}

func GetToken(usename string, opendid string) (string, error) {
	cure := time.Now()
	c := TaoClaims{
		usename,
		opendid,
		jwt.TimePrecision.String(),
		jwt.RegisteredClaims{
			Issuer:    "XianXian",
			Subject:   "View",
			Audience:  []string{"Tao"},
			ID:        opendid,
			ExpiresAt: jwt.NewNumericDate(cure.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(cure),
			NotBefore: jwt.NewNumericDate(cure),
		},
	}

	toekn := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := toekn.SignedString([]byte(GlbBaInfa.Conf.Http.Jwt))
	if err != nil {
		Logger.Infof("Get token error:%s", err.Error())
		return "", err
	}
	return ss, nil
}

func VerifyToken(tokenString string, c *gin.Context) error {
	token, err := jwt.ParseWithClaims(tokenString, &TaoClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(GlbBaInfa.Conf.Http.Jwt), nil
	})
	if err != nil {
		Logger.Infof("Parse token error:%s", err.Error())
		return err
	}

	claims, ok := token.Claims.(*TaoClaims)
	if !ok {
		Logger.Infof("Parse token error: type error")
		return errors.New("type error")
	}
	c.Set("openId", claims.OpenId)
	c.Set("Username", claims.Username)
	return nil

}
