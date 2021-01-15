package main

import (
	"net/http"
	"regexp"

	"github.com/jgolang/apirest"
	"github.com/jgolang/log"
	"github.com/jhuygens/db/users"
	"github.com/jhuygens/security"
)

func validateRequestSignUpValuesFormat(request signUpRequest) apirest.Response {
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
	validID := regexp.MustCompile(emaiFormatlLayout)
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
	user, _ := users.GetByEmail(email)
	if user != nil {
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
			Title:      errorTitle,
			Message:    "No fue posible registrar al usuario",
			ErrorCode:  "5",
			StatusCode: http.StatusInternalServerError,
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

func validateRequestResetSecretValuesFormat(request resetClientSecretRequest) apirest.Response {
	errorTitle := "Parametro inválido"
	if request.ClientID == "" {
		return apirest.Error{
			Title:   errorTitle,
			Message: "Ingresa una su client Id",
		}
	}
	if request.SecretKey == "" {
		return apirest.Error{
			Title:   errorTitle,
			Message: "Ingresa una su secret Key",
		}
	}
	if request.Password == "" {
		return apirest.Error{
			Title:   errorTitle,
			Message: "Ingresa una su contraseña",
		}
	}
	return nil
}

func resetSecret(request resetClientSecretRequest) apirest.Response {
	errorTitle := "Secret key no actualizado"
	user, err := users.GetByClientID(request.ClientID)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:      errorTitle,
			Message:    "Credenciales inválidas",
			ErrorCode:  "2",
			StatusCode: http.StatusUnauthorized,
		}
	}
	valid, err := security.ValidateSecurePassword(request.Password, user.Password)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:      errorTitle,
			Message:    "Credenciales inválidas",
			ErrorCode:  "2",
			StatusCode: http.StatusUnauthorized,
		}
	}
	if !valid {
		return apirest.Error{
			Title:      errorTitle,
			Message:    "Credenciales inválidas",
			ErrorCode:  "2",
			StatusCode: http.StatusUnauthorized,
		}
	}
	if request.SecretKey != user.CurrentSecretKey {
		return apirest.Error{
			Title:      errorTitle,
			Message:    "Credenciales inválidas",
			ErrorCode:  "2",
			StatusCode: http.StatusUnauthorized,
		}
	}

	secretKey, err := security.GenerateSecretKey(user.Email, user.ClientID, request.SecretKey, defaultSecretExpire)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:     errorTitle,
			Message:   "No fue posible generar el id del usuario",
			ErrorCode: "2",
		}
	}
	err = users.UpdateCurrentSecretKey(user.ClientID, secretKey)
	if err != nil {
		log.Error(err)
		return apirest.Error{
			Title:      errorTitle,
			Message:    "No fue posible generar el id del usuario",
			ErrorCode:  "5",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return apirest.Success{
		Title:   "¡Usuario registrado!",
		Message: "El secret actualizado correctamente",
		Data: resetClientSecretResponse{
			ClientID:  user.ClientID,
			SecretKey: secretKey,
		},
	}
}
