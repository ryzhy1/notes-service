package oauth

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
	"log/slog"
	"net/http"
	"time"
)

func AuthAPI(r *chi.Mux) {
	s := oauth.NewBearerServer(
		"yaroslav-the-best",
		time.Hour*2,
		&TestUserVerifier{},
		nil)
	r.Post("/token", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received /token request")
		s.UserCredentials(w, r)
	})
	r.Post("/auth", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received /auth request")
		s.ClientCredentials(w, r)
	})
}

type TestUserVerifier struct {
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (*TestUserVerifier) ValidateUser(username, password, scope string, r *http.Request) error {
	if username == "user1" && password == "password1" {
		return nil
	}

	return errors.New("wrong user")
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (*TestUserVerifier) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	if clientID == "abcdef" && clientSecret == "12345" {
		return nil
	}

	return errors.New("wrong client")
}

// ValidateCode validates token ID
func (*TestUserVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "", nil
}

// AddClaims provides additional claims to the token
func (*TestUserVerifier) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	claims := make(map[string]string)
	claims["customer_id"] = "1001"
	claims["customer_data"] = `{"order_date":"2016-12-14","order_id":"9999"}`
	return claims, nil
}

// AddProperties provides additional information to the token response
func (*TestUserVerifier) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	props := make(map[string]string)
	props["customer_name"] = "Gopher"
	return props, nil
}

// ValidateTokenID validates token ID
func (*TestUserVerifier) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}

// StoreTokenID saves the token id generated for the user
func (*TestUserVerifier) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
