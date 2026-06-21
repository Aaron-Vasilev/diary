package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func ParseInitData(initData string) (*User, error) {
	if err := verifyInitData(initData); err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, fmt.Errorf("no user field in initData")
	}

	var u User
	if err := json.Unmarshal([]byte(userJSON), &u); err != nil {
		return nil, err
	}

	return &u, nil
}

func verifyInitData(initData string) error {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return err
	}

	hash := values.Get("hash")
	if hash == "" {
		return fmt.Errorf("missing hash")
	}

	var pairs []string
	for k, vs := range values {
		if k == "hash" {
			continue
		}
		pairs = append(pairs, k+"="+vs[0])
	}
	sort.Strings(pairs)
	dataCheckString := strings.Join(pairs, "\n")

	secretKey := hmacSHA256([]byte("WebAppData"), []byte(os.Getenv("TELEGRAM_BOT_TOKEN")))
	computed := hex.EncodeToString(hmacSHA256(secretKey, []byte(dataCheckString)))

	if computed != hash {
		return fmt.Errorf("invalid initData hash")
	}

	return nil
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
