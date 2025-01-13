package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"starterkit/internal/config"
	"starterkit/internal/service/jwthelper"
	"strconv"
	"time"
)

//type GoogleAuthHandler struct {
//	GoogleUserStore store.GoogleUserStore
//}
//
//func NewGoogleAuthHandler(googleUserStore store.GoogleUserStore) *GoogleAuthHandler {
//	return &GoogleAuthHandler{
//		googleUserStore,
//	}
//}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (h *GoogleAuthHandler) StartGoogleOAuth(w http.ResponseWriter, r *http.Request) {
	url := config.GetAuthConfig().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *GoogleAuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := config.GetAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Use token to fetch user info
	client := config.GetAuthConfig().Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode user info
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	expirationInHoursString := os.Getenv("JWT_EXPIRATION_IN_HOURS")
	expirationInHours, err := strconv.Atoi(expirationInHoursString)

	//Expiration time for JWT and Cookie
	expirationTime := time.Now().Add(time.Duration(expirationInHours) * time.Hour)

	//Create JWT token with the user-id
	jwtToken, err := jwthelper.CreateJWT(userInfo.ID, expirationTime)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	//Create cookie for JWT
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	var newGoogleUser = store.GoogleUser{
		GoogleId:      userInfo.ID,
		Email:         userInfo.Email,
		VerifiedEmail: userInfo.VerifiedEmail,
		Name:          userInfo.Name,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
	}

	googleUser, err := h.GoogleUserStore.GetGoogleUserWhereGoogleId(newGoogleUser.GoogleId)
	if err != nil {
		var userNotFound = !errors.Is(err, sql.ErrNoRows)
		if !userNotFound {
			http.Error(w, "Failed to get google user", http.StatusInternalServerError)
		}
	}
	if googleUser == nil {
		//First time login, welcome to the platform
		err := h.GoogleUserStore.CreateGoogleUser(&newGoogleUser)
		if err != nil {
			http.Error(w, "Could not persist new user", http.StatusInternalServerError)
		}
		googleUser = &newGoogleUser
	} else {
		err = h.GoogleUserStore.UpdateGoogleUser(&newGoogleUser)
		if err != nil {
			http.Error(w, "Could not update user", http.StatusInternalServerError)
		}
	}

	// Redirect user to home
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
	http.Redirect(w, r, "/", http.StatusFound)
}
