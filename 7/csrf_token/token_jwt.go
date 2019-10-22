package main

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JwtToken struct {
	Secret []byte
}

func NewJwtToken(secret string) (*JwtToken, error) {
	return &JwtToken{Secret: []byte(secret)}, nil
}

type JwtCsrfClaims struct {
	SessionID string `json:"sid"`
	UserID    uint32 `json:"uid"`
	jwt.StandardClaims
}

func (tk *JwtToken) Create(s *Session, tokenExpTime int64) (string, error) {
	data := JwtCsrfClaims{
		SessionID: s.ID,
		UserID:    s.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tk.Secret)
}

func (tk *JwtToken) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return tk.Secret, nil
}

func (tk *JwtToken) Check(s *Session, inputToken string) (bool, error) {
	payload := &JwtCsrfClaims{}
	_, err := jwt.ParseWithClaims(inputToken, payload, tk.parseSecretGetter)
	if err != nil {
		return false, fmt.Errorf("cant parse jwt token: %v", err)
	}
	if payload.Valid() != nil {
		return false, fmt.Errorf("invalid jwt token: %v", err)
	}
	return payload.SessionID == s.ID && payload.UserID == s.UserID, nil
}

/*
func (tk *JwtToken) parseSecretGetterMultiKeys(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}

	keys := []*Key{
		&Key{Exp: 10, Secret: 1},
		&Key{Exp: 20, Secret: 2},
		&Key{Exp: 30, Secret: 3},
	}

	payload, ok := token.Claims.(*JwtCsrfClaims)
	if !ok {
		return nil, err
	}
	secret := ""
	for _, key := range keys {
		if Key.Exp > payload.Exp {
			secret = key.Secret
			break
		}
	}
	if secret == "" {
		return nil, fmt.Errrof("no secret found")
	}

	return secret, nil
}
*/
