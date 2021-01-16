package main

import (
	"net/http"
	"strings"

	"github.com/jgolang/api"
)

func generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request tokenRequest
	response := api.UnmarshalBody(&request, r)
	if response != nil {
		response.Send(w)
		return
	}

	if request.GarantType == "authorization_code" {
		auth, response := api.GetHeaderValueString("Authorization", r)
		values := strings.SplitN(auth, " ", 2)
		if len(values) != 2 || values[0] != "Basic" {
			api.Error{
				Title:      api.DefaultUnauthorizedTitle,
				Message:    api.DefaultInvalidAuthHeaderMsg,
				StatusCode: http.StatusUnauthorized,
			}.Send(w)
			return
		}
		clientID, secretKey, valid := api.ValidateBasicToken(values[1])
		if !valid {
			api.Error{
				Title:      api.DefaultUnauthorizedTitle,
				Message:    api.DefaultBasicUnauthorizedMsg,
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
	api.Error{
		Title:      api.DefaultUnauthorizedTitle,
		Message:    "garant_type no soportado",
		StatusCode: http.StatusBadRequest,
	}.Send(w)
}
