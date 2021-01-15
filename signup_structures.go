package main

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

type resetClientSecretRequest struct {
	ClientID  string `json:"client_id"`
	SecretKey string `json:"secret_key"`
	Password  string `json:"password"`
}

type resetClientSecretResponse struct {
	ClientID  string `json:"client_id"`
	SecretKey string `json:"secret_key"`
}
