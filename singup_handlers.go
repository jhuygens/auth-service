package main

import (
	"net/http"

	"github.com/jgolang/api"
)

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var request signUpRequest
	response := api.UnmarshalBody(&request, r)
	if response != nil {
		response.Write(w)
		return
	}
	response = validateRequestSignUpValuesFormat(request)
	if response != nil {
		response.Write(w)
		return
	}
	response = validateUserCreated(request.Email)
	if response != nil {
		response.Write(w)
		return
	}
	response = signUp(request)
	response.Write(w)
	return
}

func resetSecretHandler(w http.ResponseWriter, r *http.Request) {
	var request resetClientSecretRequest
	response := api.UnmarshalBody(&request, r)
	if response != nil {
		response.Write(w)
		return
	}
	response = validateRequestResetSecretValuesFormat(request)
	if response != nil {
		response.Write(w)
		return
	}
	response = resetSecret(request)
	response.Write(w)
	return
}
