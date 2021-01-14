package main

import (
	"net/http"

	"github.com/jgolang/apirest"
)

func singUpHandler(w http.ResponseWriter, r *http.Request) {
	var request signUpRequest
	response := apirest.UnmarshalBody(&request, r)
	if response != nil {
		response.Send(w)
		return
	}

	response = validateRequestValuesFormatt(request)
	if response != nil {
		response.Send(w)
		return
	}

	response = validateUserCreated(request.Email)
	if response != nil {
		response.Send(w)
		return
	}

	response = signUp(request)
	response.Send(w)
	return
}
