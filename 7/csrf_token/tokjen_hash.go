package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// fbc1fd86ab53d52c3ffeb6529aea9676e14bc52b792414c32f5612b4eb2c9745:1567618546
// JSv5M7FZ5iPHnHiLXR1QbhnMcdoY/wvEae4a76KrGBxeHruFb1S90d4GkwsoQQU4R1zqEdSa0KMGflriF2dHj5XWm4Zp6OBxLp6BJFUhqpQxEBEr5yl4sxEHadgssvVWfWtDKe0bENU=
// JSv5M7FZ5iPHnHiLXR1QbhnMcDAU/wvEae4a76KrGBxeHruFb1S90d4GkwsoQQU4R1zqEdSa0KMGflriF2dHj5XWm4Zp6OBxLp6BJFUhqpQxEBEr5yl4sxEHadgssvVWfWtDKe0bENU=
// JSv5M7FZ5iPHnHiLXR1QbhnMcNEF/wvEae4a76KrGBxeHruFb1S90d4GkwsoQQU4R1zqEdSa0KMGflriF2dHj5XWm4Zp6OBxLp6BJFUhqpQxEBEr5yl4sxEHadgssvVWfWtDKe0bENU=

type HashToken struct {
	Secret []byte
}

func NewHMACHashToken(secret string) (*HashToken, error) {
	return &HashToken{Secret: []byte(secret)}, nil
}

func (tk *HashToken) Create(s *Session, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, []byte(tk.Secret))
	data := fmt.Sprintf("%s:%d:%d", s.ID, s.UserID, tokenExpTime)
	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTime, 10)
	return token, nil
}

func (tk *HashToken) Check(s *Session, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != 2 {
		return false, fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("bad token time")
	}

	if tokenExp < time.Now().Unix() {
		return false, fmt.Errorf("token expired")
	}

	h := hmac.New(sha256.New, []byte(tk.Secret))
	data := fmt.Sprintf("%s:%d:%d", s.ID, s.UserID, tokenExp)
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, fmt.Errorf("cand hex decode token")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}
