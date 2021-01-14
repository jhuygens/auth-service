package main

import (
	"github.com/jgolang/apirest"
	"github.com/jgolang/log"
	"github.com/jhuygens/security"
)

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

	// user := users.User{
	// 	Email:            request.Email,
	// 	Password:         securePassword,
	// 	AppName:          request.AppName,
	// 	ClientID:         clientID,
	// 	CurrentSecretKey: secretKey,
	// }

	// for _, url := range request.RedirectUrls {
	// 	user.RedirectURLs = append(
	// 		user.RedirectURLs,
	// 		users.RedirectURL{
	// 			URL: url,
	// 		},
	// 	)
	// }

	// err = users.Save(user)
	// if err != nil {
	// 	log.Error(err)
	// 	return apirest.Error{
	// 		Title:     errorTitle,
	// 		Message:   "No fue posible registrar al usuario",
	// 		ErrorCode: "1",
	// 	}
	// }

	return apirest.Success{
		Title:   "¡Usuario registrado!",
		Message: "El usuario ha sido registrado correctamente",
		Data: signUpResponse{
			ClientID:  clientID,
			SecretKey: secretKey,
		},
	}
}
