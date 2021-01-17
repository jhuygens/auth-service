package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jgolang/api"
	"github.com/jgolang/config"
	"github.com/jgolang/log"
	_ "github.com/jhuygens/db-mongodb/users"
	"github.com/jhuygens/security"
)

const emailFormatLayout = `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`

var (
	defaultSecretExpire = config.GetInt("general.default_secret_expire")
	defaultTokenExpire  = config.GetInt("general.default_token_expire")
)

func init() {
	// Api package custom config
	api.CustomTokenValidatorFunc = security.ValidateAccessTokenFunc
	api.Print = log.Printf
	api.PrintError = log.Error
	api.Fatal = log.Fatal
}

func main() {
	router := mux.NewRouter()
	port := config.GetInt("services.auth.port")
	api.CustomTokenValidatorFunc = security.ValidateAccessTokenFunc
	api.ValidateBasicAuthCredentialsFunc = validateBasicAuthCredentials
	noAuthMiddlewares := api.MiddlewaresChain(api.ContentExtractor)
	tokenAuthMiddlewares := api.MiddlewaresChain(api.ContentExtractor, api.CustomToken)

	router.HandleFunc("/v1/signup", noAuthMiddlewares(signUpHandler)).Methods(http.MethodPost)
	router.HandleFunc("/v1/reset_secret", tokenAuthMiddlewares(resetSecretHandler)).Methods(http.MethodPost)
	router.HandleFunc("/v1/token", noAuthMiddlewares(generateTokenHandler)).Methods(http.MethodPost)

	log.Info("Starting server, listen on port: ", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), router); err != nil {
		log.Panic(err)
	}
}
