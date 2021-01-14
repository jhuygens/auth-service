package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jgolang/apirest"
	"github.com/jgolang/config"
	"github.com/jgolang/log"
)

var (
	defaultSecretExpire = config.GetInt("general.default_secret_expire")
	// EventID doc ...
	additionalInfo = `{\"service\":\"auth-service\",\"event_id\":\"\"}`
)

func init() {
	log.OverrideConfig(log.LstdDevFlags, log.LstdProdFlags|log.Linfo, &additionalInfo)
}

func main() {
	router := mux.NewRouter()
	port := config.GetInt("services.auth.port")
	middlewares := apirest.MiddlewaresChain(apirest.ContentExtractor)

	router.HandleFunc("/v1/signup", middlewares(singUpHandler)).Methods(http.MethodPost)

	log.Info("Starting server, lisen on port: ", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), router); err != nil {
		log.Panic(err)
	}
}

func setEventID(eventID string) {
	additionalInfo = `{\"project\":\"jhuygens\",\"service\":\"auth-service\",\"event_id\":\"` + eventID + `\"}`
}
