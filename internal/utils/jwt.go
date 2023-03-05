package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
	jwt.StandardClaims
}

var SecretKey = ""

func GenerateToken(id uint, fullname string) (string, error) {
	claims := &Claims{
		ID:       id,
		Fullname: fullname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(48 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	webtoken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return webtoken, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func DecodeToken(tokenString string) (*Claims, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Type of token.Claims: %T\n", token.Claims)
	claims, isOk := token.Claims.(jwt.MapClaims)
	if isOk && token.Valid {
		authData, _ := json.Marshal(claims)
		var res Claims
		json.Unmarshal(authData, &res)

		return &res, nil
	}

	return nil, fmt.Errorf("invalid token")
}
