package main

import (
	"net/http"

	"github.com/jgolang/apirest"
)

func singUpHandler(w http.ResponseWriter, r *http.Request) {
	eventID, response := apirest.GetHeaderValueString("EventID", r)
	if response != nil {
		response.Send(w)
		return
	}

	go setEventID(eventID)

	var requestObject signUpRequest
	response = apirest.UnmarshalBody(&requestObject, r)
	if response != nil {
		response.Send(w)
		return
	}

	response = requestObject.ValidateRequestValuesFormatt()
	if response != nil {
		response.Send(w)
		return
	}

	response = signUp(requestObject)
	response.Send(w)
	return
}
