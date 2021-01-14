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
	defaultTokenExpire  = config.GetInt("general.default_token_expire")
)

func main() {
	router := mux.NewRouter()
	port := config.GetInt("services.auth.port")
	middlewares := apirest.MiddlewaresChain(apirest.ContentExtractor)

	router.HandleFunc("/v1/signup", middlewares(singUpHandler)).Methods(http.MethodPost)
	router.HandleFunc("/v1/reset_secret", middlewares(singUpHandler)).Methods(http.MethodPost)
	router.HandleFunc("/v1/token", middlewares(singUpHandler)).Methods(http.MethodPost)

	log.Info("Starting server, lisen on port: ", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), router); err != nil {
		log.Panic(err)
	}
}
