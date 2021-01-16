package main

import (
	"net/http"

	"github.com/jhuygens/db/users"

	"github.com/jgolang/api"
	"github.com/jgolang/log"
	"github.com/jhuygens/security"
)

func generateToken(clientID, secretID string) api.Response {
	errorTitle := "No autorizado"
	token, err := security.GenerateAccessToken(clientID, secretID, defaultTokenExpire)
	if err != nil {
		return api.Error{
			Title:     errorTitle,
			Message:   "No se ha podido general el token",
			ErrorCode: "2",
		}
	}
	refreshToken, err := security.GenerateAccessToken(clientID, token, defaultTokenExpire)
	if err != nil {
		return api.Error{
			Title:     errorTitle,
			Message:   "No se ha podido generar el token",
			ErrorCode: "2",
		}
	}
	err = users.UpdateToken(clientID, token, refreshToken)
	if err != nil {
		log.Error(err)
		return api.Error{
			Title:      errorTitle,
			Message:    "No es posible actualizar el token",
			ErrorCode:  "5",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return api.Success{
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

func refreshToken(request tokenRequest) api.Response {
	errorTitle := "No autorizado"
	token, err := security.GenerateAccessToken(request.ClientID, request.RefreshToken, defaultTokenExpire)
	if err != nil {
		return api.Error{
			Title:     errorTitle,
			Message:   "No se ha podido generar el token",
			ErrorCode: "2",
		}
	}
	refreshToken, err := security.GenerateAccessToken(request.ClientID, token, defaultTokenExpire)
	if err != nil {
		return api.Error{
			Title:     errorTitle,
			Message:   "No se ha podido generar el token",
			ErrorCode: "2",
		}
	}
	err = users.UpdateToken(request.ClientID, token, refreshToken)
	if err != nil {
		log.Error(err)
		return api.Error{
			Title:      errorTitle,
			Message:    "No es posible actualizar el token",
			ErrorCode:  "5",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return api.Success{
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

func validateRefreshToken(clientID, token string) api.Response {
	errorTitle := "No autorizado"
	if token == "" {
		return api.Error{
			Title:      api.DefaultUnauthorizedTitle,
			Message:    "Token de refresco inválido",
			StatusCode: http.StatusUnauthorized,
		}
	}
	user, err := users.GetByClientID(clientID)
	if err != nil {
		log.Error(err)
		return api.Error{
			Title:      errorTitle,
			Message:    "No es posible actualizar el token",
			ErrorCode:  "5",
			StatusCode: http.StatusUnauthorized,
		}
	}
	if token != user.RefreshToken {
		return api.Error{
			Title:      errorTitle,
			Message:    "Token de refresco inválido",
			StatusCode: http.StatusUnauthorized,
		}
	}
	return nil
}
