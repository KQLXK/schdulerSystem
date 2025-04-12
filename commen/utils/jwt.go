package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"schedule/commen/config"
	"schedule/models"
	"time"
)

type Jwtinfo struct {
	Identity string
	UserID   int64
}

type claims struct {
	Jwtinfo
	jwt.RegisteredClaims
}

var (
	SigningMethod   = jwt.SigningMethodHS256
	TokenExpiredErr = errors.New("token expired")
)

func (i *Jwtinfo) CreateToken() (string, error) {
	accesstoken, err := i.CreateAccessToken()
	if err != nil {
		return "", err
	}
	refreshtoken, err := i.CreateRefreshToken()
	if err != nil {
		return "", err
	}
	err = models.NewTokenDao().Setkey(accesstoken, refreshtoken)
	if err != nil {
		return "", err
	}
	return accesstoken, nil
}

func (info *Jwtinfo) CreateAccessToken() (string, error) {
	Myclaims := claims{
		Jwtinfo: Jwtinfo{
			Identity: info.Identity,
			UserID:   info.UserID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(SigningMethod, Myclaims)
	acctoken, err := token.SignedString([]byte(config.GetConfig().JWT.SecretKey))
	if err != nil {
		log.Println("create accesstoken failed, err:", err)
		return "", err
	}
	return acctoken, nil
}

func (info *Jwtinfo) CreateRefreshToken() (string, error) {
	Myclaims := claims{
		Jwtinfo: Jwtinfo{
			Identity: info.Identity,
			UserID:   info.UserID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 10)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(SigningMethod, Myclaims)
	reftoken, err := token.SignedString([]byte(config.GetConfig().JWT.SecretKey))
	if err != nil {
		log.Println("create refreshtoken failed, err:", err)
		return "", err
	}
	return reftoken, nil
}

func GetInfoFromToken(Signstring string) (*Jwtinfo, error) {
	token, err := jwt.ParseWithClaims(Signstring, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWT.SecretKey), nil
	})
	if err != nil {
		log.Println("parse token failed, err:", err)
		if err == jwt.ErrTokenExpired {
			return nil, TokenExpiredErr
		}
		return nil, err
	}
	if info, ok := token.Claims.(*claims); ok && token.Valid {
		return &info.Jwtinfo, nil
	}
	return nil, errors.New("invalid token")
}

//func Test() {
//	info := &Jwtinfo{
//		Identity: "teacher",
//		UserID:   10,
//	}
//	token, _ := CreateAccessToken(info)
//	fmt.Println(token)
//	info, _ = GetInfoFromToken(token)
//	fmt.Println(info)
//}
