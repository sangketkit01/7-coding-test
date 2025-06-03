package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Maker interface {
	CreateToken(objectId primitive.ObjectID, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

type JWTMaker struct {
	secretKey string
}

func NewMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, errors.New("invalid secret key length, must be at least 32 characters")
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func (m *JWTMaker) CreateToken(objectId primitive.ObjectID, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(objectId, duration)
	if err != nil {
		return "", nil,err
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwt.SignedString([]byte(m.secretKey))
	
	return token, payload, err
}

func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Invalid token.")
		}

		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, fmt.Errorf("token is invalid")) {
			return nil, fmt.Errorf("token is invalid")
		}

		return nil, fmt.Errorf("token is invalid")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf("token is invalid")
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
