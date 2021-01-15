package main

import (
	"net/http"
	"strings"

	"github.com/jgolang/apirest"
)

func generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request tokenRequest
	response := apirest.UnmarshalBody(&request, r)
	if response != nil {
		response.Send(w)
		return
	}

	if request.GarantType == "authorization_code" {
		auth, response := apirest.GetHeaderValueString("Authorization", r)
		values := strings.SplitN(auth, " ", 2)
		if len(values) != 2 || values[0] != "Basic" {
			apirest.Error{
				Title:      apirest.DefaultUnauthorizedTitle,
				Message:    apirest.DefaultInvalidAuthHeaderMsg,
				StatusCode: http.StatusUnauthorized,
			}.Send(w)
			return
		}
		clientID, secretKey, valid := apirest.ValidateBasicToken(values[1])
		if !valid {
			apirest.Error{
				Title:      apirest.DefaultUnauthorizedTitle,
				Message:    apirest.DefaultBasicUnauthorizedMsg,
				StatusCode: http.StatusUnauthorized,
			}.Send(w)
			return
		}
		response = generateToken(clientID, secretKey)
		response.Send(w)
		return
	}

	if request.GarantType == "refresh_token" {
		response := validateRefreshToken(request.ClientID, request.RefreshToken)
		if response != nil {
			response.Send(w)
			return
		}
		response = refreshToken(request)
		response.Send(w)
		return
	}

}
