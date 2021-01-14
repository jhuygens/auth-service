package main

import (
	"regexp"

	"github.com/jgolang/apirest"
	"github.com/jgolang/log"
	"github.com/jhuygens/db/users"
	"github.com/jhuygens/security"
)

func validateRequestValuesFormat(request signUpRequest) apirest.Response {
	errorTitle := "Parametro inválido"
	if !validateEmailFormat(request.Email) {
		return apirest.Error{
			Title:   errorTitle,
			Message: "Ingresa una dirección de correo electrónico válida",
		}
	}
	if !validatePasswordFormat(request.Password) {
		return apirest.Error{
			Title:   errorTitle,
			Message: "Tu contraseña debe tener al menos 6 caracteres",
		}
	}
	if request.AppName == "" {
		return apirest.Error{
			Title:   errorTitle,
			Message: "Debes ingresar el nombre de tu applicación",
		}
	}
	if len(request.AppName) < 6 {
		return apirest.Error{
			Title:   errorTitle,
			Message: "El nombre de tu applicación debe tener al menos 6 caracteres",
		}
	}
	return nil
}

func validateEmailFormat(email string) bool {
	emailLayout := `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	validID := regexp.MustCompile(emailLayout)
	return validID.MatchString(email)
}

func validatePasswordFormat(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

func validateUserCreated(email string) apirest.Response {
	errorTitle := "Usuario no registrado"
	user, err := users.GetByEmail(email)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No es posible registrar al usuario en estos momentos, por favor intenta más tarde",
			ErrorCode: "2",
		}
	}
	if user.ID != 0 {
		return apirest.Error{
			Title:   errorTitle,
			Message: "Ya existe un usuario registrado con ese correo electronico",
		}
	}
	return nil
}

func signUp(request signUpRequest) apirest.Response {
	errorTitle := "Usuario no registrado"
	securePassword, err := security.GenerateSecurePassword(request.Password)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No se ha podido generar una contraseña para el usuario",
			ErrorCode: "2",
		}
	}
	clientID := security.CreateNewClientID(request.Email)
	secretKey, err := security.GenerateSecretKey(request.Email, clientID, securePassword, defaultSecretExpire)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No fue posible generar el id del usuario",
			ErrorCode: "2",
		}
	}
	user := users.User{
		Email:            request.Email,
		Password:         securePassword,
		AppName:          request.AppName,
		ClientID:         clientID,
		CurrentSecretKey: secretKey,
	}
	for _, url := range request.RedirectUrls {
		user.RedirectURLs = append(
			user.RedirectURLs,
			users.RedirectURL{
				URL: url,
			},
		)
	}
	err = users.Save(user)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No fue posible registrar al usuario",
			ErrorCode: "1",
		}
	}
	return apirest.Success{
		Title:   "¡Usuario registrado!",
		Message: "El usuario ha sido registrado correctamente",
		Data: signUpResponse{
			ClientID:  clientID,
			SecretKey: secretKey,
		},
	}
}
