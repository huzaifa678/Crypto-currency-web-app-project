package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
		usePKCE			  bool
) 

func InitGoogleOAuth(config utils.Config) {
	usePKCE = config.Environment == "enterprise"
	googleClientID = config.GoogleClientID 

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
    redirectTo := r.URL.Query().Get("redirect_to")
    if redirectTo == "" {
        redirectTo = "/dashboard" 
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "oauth_redirect",
        Value:    redirectTo,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, 
        SameSite: http.SameSiteLaxMode,
    })

    stateBytes := make([]byte, 16)
    _, _ = rand.Read(stateBytes)
    state := base64.RawURLEncoding.EncodeToString(stateBytes)

    http.SetCookie(w, &http.Cookie{
        Name:     "oauth_state",
        Value:    state,
        Path:     "/",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
    })

    var url string
    if usePKCE {
        verifier, err := generateCodeVerifier()
        if err != nil {
            http.Error(w, "failed to generate verifier", http.StatusInternalServerError)
            return
        }

        challenge := generateCodeChallenge(verifier)
        http.SetCookie(w, &http.Cookie{
            Name:     "pkce_verifier",
            Value:    verifier,
            Path:     "/",
            HttpOnly: true,
            Secure:   false,
            SameSite: http.SameSiteLaxMode,
        })

        url = googleOAuthConfig.AuthCodeURL(
            state,
            oauth2.AccessTypeOffline,
            oauth2.SetAuthURLParam("code_challenge", challenge),
            oauth2.SetAuthURLParam("code_challenge_method", "S256"),
        )
    } else {
        url = googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
    }

    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "missing code", http.StatusBadRequest)
        return
    }

    stateFromGoogle := r.URL.Query().Get("state")
    stateCookie, err := r.Cookie("oauth_state")
    if err != nil || stateFromGoogle != stateCookie.Value {
        http.Error(w, "invalid oauth state", http.StatusBadRequest)
        return
    }

    ctx := context.Background()
    var token *oauth2.Token

    if usePKCE {
        pkceCookie, err := r.Cookie("pkce_verifier")
        if err != nil {
            http.Error(w, "missing pkce verifier", http.StatusBadRequest)
            return
        }

        token, err = googleOAuthConfig.Exchange(
            ctx,
            code,
            oauth2.SetAuthURLParam("code_verifier", pkceCookie.Value),
        )
    } else {
        token, err = googleOAuthConfig.Exchange(ctx, code)
    }

    if err != nil {
        http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    clearCookie(w, "pkce_verifier")
    clearCookie(w, "oauth_state")

    rawIDToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "no id_token in response", http.StatusInternalServerError)
        return
    }

    payload, err := idtoken.Validate(ctx, rawIDToken, googleClientID)
    if err != nil {
        http.Error(w, "invalid ID token: "+err.Error(), http.StatusUnauthorized)
        return
    }

    emailVal, _ := payload.Claims["email"].(string)

    redirectCookie, err := r.Cookie("oauth_redirect")
    redirectURL := "http://localhost:5173/dashboard" 
    if err == nil && redirectCookie.Value != "" {
        redirectURL = redirectCookie.Value
    }

    redirectWithToken := fmt.Sprintf("%s?token=%s&email=%s", redirectURL, rawIDToken, emailVal)
    http.Redirect(w, r, redirectWithToken, http.StatusSeeOther)
}

func VerifyGoogleIDToken(ctx context.Context, rawIDToken string, clientID string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(ctx, rawIDToken, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid google id_token: %w", err)
	}
	return payload, nil
}

// helper funtion to clear cookies
func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
