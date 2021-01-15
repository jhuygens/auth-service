package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jgolang/apirest"
	"github.com/jgolang/config"
	"github.com/jgolang/log"
)

const emaiFormatlLayout = `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`

var (
	defaultSecretExpire = config.GetInt("general.default_secret_expire")
	defaultTokenExpire  = config.GetInt("general.default_token_expire")
)

func main() {
	router := mux.NewRouter()
	port := config.GetInt("services.auth.port")
	apirest.CustomTokenValidatorFunc = validateAccessTokenFunc
	apirest.ValidateBasicAuthCredentialsFunc = validateBasicAuthCredentials
	noAuthMiddlewares := apirest.MiddlewaresChain(apirest.ContentExtractor)
	tokenAuthMiddlewares := apirest.MiddlewaresChain(apirest.ContentExtractor, apirest.CustomToken)

	router.HandleFunc("/v1/signup", noAuthMiddlewares(signUpHandler)).Methods(http.MethodPost)
	router.HandleFunc("/v1/reset_secret", tokenAuthMiddlewares(resetSecretHandler)).Methods(http.MethodPost)
	router.HandleFunc("/v1/token", noAuthMiddlewares(generateTokenHandler)).Methods(http.MethodPost)

	log.Info("Starting server, lisen on port: ", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), router); err != nil {
		log.Panic(err)
	}
}
