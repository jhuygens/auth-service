package main

import (
	"encoding/json"

	"github.com/jgolang/log"
	"github.com/jhuygens/db/users"
	"github.com/jhuygens/security"
)

func validateAccessTokenFunc(token string) (json.RawMessage, bool) {
	tokenData, valid, err := security.ValidateAccessToken(token)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	return []byte(tokenData), valid
}

func validateBasicAuthCredentials(clientID, secretKey string) bool {
	user, err := users.GetByClientID(clientID)
	if err != nil {
		log.Error(err)
		return false
	}
	if secretKey != user.CurrentSecretKey {
		return false
	}
	return true
}
