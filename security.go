package main

import (
	"github.com/jgolang/log"
	"github.com/jhuygens/db/users"
)

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
