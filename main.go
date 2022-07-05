package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	muxserver "github.com/asyrawi/gojwt/mux_server"
	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	Group    string `json:"group"`
}

var APPLICATION_NAME = "SIMPLE JWT APP"
var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("secret")

// main function  
func main() {

	// Create Server
	mux := new(muxserver.CustomerServerMux)
	mux.RegisterMiddleware(JwtMiddleware)
	mux.RegisterMiddleware(LoggerMidlleware)

	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("/login", LoginHandler)

	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":8000"
	// Running Server
	fmt.Println("Running Server")
	server.ListenAndServe()

}

// RootHandler function  
func RootHandler(w http.ResponseWriter, r *http.Request) {

	result := r.Context().Value("user").(jwt.MapClaims)
	response, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprint(w, string(response))
}

// LoginHandler function  
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Harus Method Post Boss ku", http.StatusBadRequest)
		return
	}

	username, password, ok := r.BasicAuth()

	if !ok {
		http.Error(w, "Hmm Something Error", http.StatusBadRequest)
	}

	stateLogin := FakeLogin(username, password)

	if !stateLogin {
		http.Error(w, "Hmm Gagal login", http.StatusUnauthorized)
	}

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		Username: username,
		Email:    "hasyrawi@gmail.com",
		Group:    "user-1",
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return
	}

	jsonResponse := map[string]string{
		"token": signedToken,
	}

	result, err := json.Marshal(jsonResponse)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprint(w, string(result))

}

// Every Connection Will Be Interecept By Jwt
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Hanan", "jwt-request")

		// Unprotect Url
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")

		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := ParseJwt(tokenString)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(context.Background(), "user", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func ParseJwt(token string) (*jwt.Token, error) {

	tokenResult, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing Method Invalid")
		} else if method != JWT_SIGNING_METHOD {
			return nil, fmt.Errorf("Invalid Method")
		}
		return JWT_SIGNING_METHOD, nil
	})

	return tokenResult, nil
}

func FakeLogin(username, password string) bool {
	if username == "hanan" && password == "password" {
		return true
	}
	return false
}

func LoggerMidlleware(next http.Handler) http.Handler {
	// Implement the signature
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
