package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go-hacktiv8-2/config"

	"github.com/golang-jwt/jwt/v4"
)

type AuthDetails struct {
	AuthUUID string
	UserID   uint64
}

var cfg = config.LoadConfig()
var SECRET_KEY = []byte(cfg.ApiSecret)

func CreateToken(authD AuthDetails) (string, error) {
	claim := jwt.MapClaims{}
	claim["authorized"] = true
	claim["auth_uuid"] = authD.AuthUUID
	claim["user_id"] = authD.UserID
	claim["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return err
	}
	return nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}

	bearToken := r.Header.Get("Authorization")
	bearer := strings.HasPrefix(bearToken, "Bearer")
	if bearer {
		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			return strArr[1]
		}
		return ""
	}

	return ""
}

func ExtractTokenAuth(r *http.Request) (*AuthDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		authUuid, ok := claims["auth_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AuthDetails{
			AuthUUID: authUuid,
			UserID:   userId,
		}, nil
	}
	return nil, err
}
