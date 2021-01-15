package main

import (
	"net/http"

	"github.com/jhuygens/db/users"

	"github.com/jgolang/apirest"
	"github.com/jgolang/log"
	"github.com/jhuygens/security"
)

func generateToken(clientID, secretID string) apirest.Response {
	errorTitle := "No autorizado"
	token, err := security.GenerateAccessToken(clientID, secretID, defaultTokenExpire)
	if err != nil {
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No se ha podido general el token",
			ErrorCode: "2",
		}
	}
	refreshToken, err := security.GenerateAccessToken(clientID, token, defaultTokenExpire)
	if err != nil {
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No se ha podido generar el token",
			ErrorCode: "2",
		}
	}
	err = users.UpdateToken(clientID, token, refreshToken)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:      errorTitle,
			Message:    "No es posible actualizar el token",
			ErrorCode:  "5",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return apirest.Success{
		Title:   "Token autorizado!",
		Message: "Se ha generado correctamente el token",
		Data: tokenResponse{
			AccessToken:  token,
			TokenType:    "Bearer",
			ExpiresIn:    defaultTokenExpire,
			RefreshToken: refreshToken,
		},
	}
}

func refreshToken(request tokenRequest) apirest.Response {
	errorTitle := "No autorizado"
	token, err := security.GenerateAccessToken(request.ClientID, request.RefreshToken, defaultTokenExpire)
	if err != nil {
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No se ha podido generar el token",
			ErrorCode: "2",
		}
	}
	refreshToken, err := security.GenerateAccessToken(request.ClientID, token, defaultTokenExpire)
	if err != nil {
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No se ha podido generar el token",
			ErrorCode: "2",
		}
	}
	err = users.UpdateToken(request.ClientID, token, refreshToken)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:      errorTitle,
			Message:    "No es posible actualizar el token",
			ErrorCode:  "5",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return apirest.Success{
		Title:   "Token autorizado!",
		Message: "Se ha refrescado correctamente el token",
		Data: tokenResponse{
			AccessToken:  token,
			TokenType:    "Bearer",
			ExpiresIn:    defaultTokenExpire,
			RefreshToken: refreshToken,
		},
	}
}

func validateRefreshToken(clientID, token string) apirest.Response {
	errorTitle := "No autorizado"
	if token == "" {
		return apirest.Error{
			Title:      apirest.DefaultUnauthorizedTitle,
			Message:    "Token de refresco inválido",
			StatusCode: http.StatusUnauthorized,
		}
	}
	user, err := users.GetByClientID(clientID)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:      errorTitle,
			Message:    "No es posible actualizar el token",
			ErrorCode:  "5",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if token != user.RefreshToken {
		return apirest.Error{
			Title:      errorTitle,
			Message:    "Token de refresco inválido",
			StatusCode: http.StatusUnauthorized,
		}
	}
	return nil
}
