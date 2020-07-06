package utils

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	// "github.com/unknwon/com"
	"github.com/xifengzhu/eshop/initializers/setting"
)

var mySigningKey = setting.JwtSecret

func Encode(params map[string]interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(720)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["issu"] = "eshop"
	for key, value := range params {
		claims[key] = value
	}
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(mySigningKey))
	return tokenString
}

func Decode(token string) (map[string]interface{}, error) {
	parseAuth, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})
	if err != nil {
		return nil, errors.New("TOKEN ERROR!")
	}
	//将token中的内容存入parmMap
	claim := parseAuth.Claims.(jwt.MapClaims)
	var parmMap map[string]interface{}
	parmMap = make(map[string]interface{})
	for key, val := range claim {
		parmMap[key] = val
	}
	log.Print("Decode Info: ", parmMap)
	return parmMap, nil
}
