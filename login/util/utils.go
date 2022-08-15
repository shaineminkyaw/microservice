package util

import (
	"crypto/rsa"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/shaineminkyaw/microservice/login/ds"
	"github.com/shaineminkyaw/microservice/user/config"
	"github.com/shaineminkyaw/microservice/user/model"
)

type AccessTokenCustomClaim struct {
	UserID uint64
	jwt.StandardClaims
}

type refreshTokenClaim struct {
	UserID uint64
	jwt.StandardClaims
}
type RefreshTokenCustomClaim struct {
	UID         uint64
	TokenID     uuid.UUID
	TokenString string
	Expire      time.Duration
}

type AccessToken struct {
	SS string
}
type RefreshToken struct {
	Uid         uint64
	TokenID     string
	TokenString string
	Expire      time.Duration
}

type TokenPair struct {
	AccessToken
	RefreshToken
}

func GetAccessToken(userID uint64, key *rsa.PrivateKey) (string, error) {
	//
	unixTime := time.Now()
	expTime := unixTime.Add(60 * 15) // 15min

	claims := &AccessTokenCustomClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime.Unix(),
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenData, err := token.SignedString(key)
	if err != nil {
		log.Println(err.Error())
	}
	return tokenData, err
}

func ValidateAccessToken(token string, key *rsa.PublicKey) (*AccessTokenCustomClaim, error) {
	//
	claim := AccessTokenCustomClaim{}
	tokenData, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		log.Println(err.Error())
	}
	if !tokenData.Valid {
		log.Println(err.Error())
	}
	data, ok := tokenData.Claims.(*AccessTokenCustomClaim)
	if !ok {
		log.Println(err.Error())
	}
	return data, nil
}

func GetRefreshToken(userId uint64, key *rsa.PrivateKey) (*RefreshTokenCustomClaim, error) {
	//
	unixTime := time.Now()
	expTime := unixTime.Add(60 * 60 * 24 * 7) // 1week
	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Println(err.Error())
	}

	claim := &refreshTokenClaim{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime.Unix(),
			ExpiresAt: expTime.Unix(),
			Id:        tokenID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	tokenData, err := token.SignedString(key)
	if err != nil {
		log.Println(err.Error())
	}

	return &RefreshTokenCustomClaim{
		UID:         userId,
		TokenID:     tokenID,
		TokenString: tokenData,
		Expire:      expTime.Sub(unixTime),
	}, nil
}

func ValidRefreshToken(token string, key *rsa.PublicKey) (*RefreshTokenCustomClaim, error) {
	//
	claim := &refreshTokenClaim{}
	tokenData, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		log.Println(err.Error())
	}
	if !tokenData.Valid {
		log.Println(err.Error())
	}
	data, ok := tokenData.Claims.(*refreshTokenClaim)
	if !ok {
		log.Println(err.Error())
	}
	uuidData, _ := uuid.Parse(data.Id)
	return &RefreshTokenCustomClaim{
		UID:         data.UserID,
		TokenID:     uuidData,
		TokenString: token,
		Expire:      time.Duration(data.ExpiresAt),
	}, nil
}

//Delete Token

func DeleteRefreshToken(userId uint64) error {
	//
	token := &model.UserToken{}
	err := ds.Login_DB.Model(&model.UserToken{}).Where("uid = ?", userId).Delete(&token).Error
	if err != nil {
		return err
	}
	return nil
}

///Generate New Token Pair

func NewTokenPair(user *model.User, prvToken string) (*TokenPair, error) {
	//
	if len(prvToken) > 0 {
		err := DeleteRefreshToken(user.ID)
		if err != nil {
			log.Print(err.Error())
		}
	}

	conf := config.Init()
	accessToken, err := GetAccessToken(user.ID, conf.Private)
	if err != nil {
		log.Print(err.Error())
	}
	refreshToken, err := GetRefreshToken(user.ID, conf.Private)
	if err != nil {
		log.Print(err.Error())
	}
	rToken, err := ValidRefreshToken(refreshToken.TokenString, conf.Public)
	if err != nil {
		log.Print(err.Error())
	}
	token := &model.UserToken{
		Uid:        user.ID,
		TokenID:    rToken.TokenID.String(),
		Token:      rToken.TokenString,
		ExpireTime: time.Now().Add(rToken.Expire),
	}

	err = ds.Login_DB.Model(&model.UserToken{}).Create(&token).Error
	if err != nil {
		log.Print(err.Error())
	}

	return &TokenPair{
		AccessToken: AccessToken{
			accessToken,
		},
		RefreshToken: RefreshToken{
			Uid:         token.Uid,
			TokenID:     token.TokenID,
			TokenString: token.Token,
			Expire:      time.Since(token.ExpireTime),
		},
	}, nil
}
