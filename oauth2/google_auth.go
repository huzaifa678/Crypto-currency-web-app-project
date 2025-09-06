package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var (
		googleOAuthConfig *oauth2.Config
		googleClientID    string
) 

func InitGoogleOAuth(config utils.Config) {
    googleOAuthConfig = &oauth2.Config{
        ClientID:     config.GoogleClientID,
        ClientSecret: config.GoogleClientSecret,
        RedirectURL:  config.GoogleRedirectURL,
        Scopes: []string{
            "openid", 
            "email",
            "profile",
        },
        Endpoint: google.Endpoint,
    }
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("state-random", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token in response", http.StatusInternalServerError)
		return
	}

	payload, err := idtoken.Validate(context.Background(), rawIDToken, googleClientID)
	if err != nil {
		http.Error(w, "invalid ID token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	emailVal, ok := payload.Claims["email"].(string)
	if !ok {
    	emailVal = ""
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"id_token": rawIDToken,
		"email_validation":    emailVal,
	})
}

func VerifyGoogleIDToken(ctx context.Context, rawIDToken string, clientID string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(ctx, rawIDToken, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid google id_token: %w", err)
	}
	return payload, nil
}
