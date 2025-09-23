package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"server/utils" // Assuming this package provides GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// UserData represents the user information returned from Google's API.
type UserData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// googleOAuthConfig holds the Google OAuth2 configuration.
var googleOAuthConfig *oauth2.Config

// InitGoogleOAuth initializes the Google OAuth2 configuration.
func InitGoogleOAuth() {
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/google/callback",
		ClientID:     utils.GOOGLE_CLIENT_ID,
		ClientSecret: utils.GOOGLE_CLIENT_SECRET,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleLogin handles the redirect to Google's login page.
// The handler must have the signature http.HandlerFunc: func(w http.ResponseWriter, r *http.Request)
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// The state string is a CSRF token. It should be unique per user.
	// In a real app, generate and store it securely (e.g., in a cookie or session).
	state := "randomstate"
	url := googleOAuthConfig.AuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the callback from Google after a successful login.
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// The 'state' and 'code' are retrieved from the request's query parameters.
	state := r.URL.Query().Get("state")
	if state != "randomstate" {
		http.Error(w, "States don't match.", http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not provided.", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an OAuth2 token.
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code-Token Exchange Failed: %v", err)
		http.Error(w, "Failed to exchange authorization code for a token.", http.StatusInternalServerError)
		return
	}

	// Use the token to get user information from Google's API.
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("User Data Fetch Failed: %v", err)
		http.Error(w, "Failed to fetch user data.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Google API returned status code %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	userDataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read user data.", http.StatusInternalServerError)
		return
	}

	var userData UserData
	if err := json.Unmarshal(userDataBytes, &userData); err != nil {
		log.Printf("JSON Unmarshaling Failed: %v", err)
		http.Error(w, "Failed to parse user data.", http.StatusInternalServerError)
		return
	}

	// For demonstration, we'll return the user data as a JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(userDataBytes); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
