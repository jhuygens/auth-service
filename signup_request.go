package main

import (
	"regexp"

	"github.com/jgolang/apirest"
)

type signUpRequest struct {
	Email        string   `json:"email"`
	Password     string   `json:"password"`
	AppName      string   `json:"app_name"`
	RedirectUrls []string `json:"redirect_urls"`
}

type signUpResponse struct {
	ClientID  string `json:"client_id"`
	SecretKey string `json:"secret_key"`
}

func (request signUpRequest) ValidateRequestValuesFormatt() apirest.Response {
	errorTitle := "Parametro inválido"
	errorCode := ""
	if !validateEmailFormat(request.Email) {
		return apirest.Error{
			Title:     errorTitle,
			Message:   "Ingresa una dirección de correo electrónico válida",
			ErrorCode: errorCode,
		}
	}
	if !validatePasswordFormat(request.Password) {
		return apirest.Error{
			Title:     errorTitle,
			Message:   "Tu contraseña debe tener al menos 6 caracteres",
			ErrorCode: errorCode,
		}
	}
	if request.AppName == "" {
		return apirest.Error{
			Title:     errorTitle,
			Message:   "Debes ingresar el nombre de tu applicación",
			ErrorCode: errorCode,
		}
	}
	if len(request.AppName) < 6 {
		return apirest.Error{
			Title:     errorTitle,
			Message:   "El nombre de tu applicación debe tener al menos 6 caracteres",
			ErrorCode: errorCode,
		}
	}
	// if len(request.RedirectUrls) == 0 {
	// 	return apirest.Error{
	// 		Title:   errorTitle,
	// 		Message: "Debes ingresar al menos una url",
	// 		ErrorCode: errorCode,
	// 	}
	// }
	return nil
}

func validateEmailFormat(email string) bool {
	var validID = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	return validID.MatchString(email)
}

func validatePasswordFormat(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}
