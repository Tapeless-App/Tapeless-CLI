package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/env"
	versionService "tapeless.app/tapeless-cli/services/version"
	"tapeless.app/tapeless-cli/util"
)

func EnsureValidSession() (string, error) {

	token := viper.GetString("token")

	if token == "" {
		fmt.Println("No token found. Please log in.")
		return FetchNewToken()
	}

	isExpired, err := util.IsJWTExpired(token)

	if err != nil {
		fmt.Println("Error verifying access token:", err)
		return FetchNewToken()
	}

	if isExpired {
		fmt.Println("JWT token expired.")
		return FetchNewToken()
	}

	return token, nil

}

func FetchNewToken() (string, error) {

	versionService.CheckLatestVersion()

	loginURL := env.WebURL + "/cli/login"
	fmt.Println("Opening browser to log in...")
	fmt.Println("If the browser does not open automatically, go to:", loginURL)

	time.Sleep(1 * time.Second)

	// Open browser for user to log in
	err := util.OpenBrowser(loginURL)
	if err != nil {
		return "", err
	}

	// Optional: Wait for JWT (via a callback server or polling mechanism)
	jwt, err := waitForJWT()
	if err != nil {
		return "", err
	}

	// Store the JWT in your config
	viper.Set("token", jwt)
	viper.WriteConfig()

	return jwt, nil

}

// waitForJWT starts a local server and waits for the JWT to be sent back from the web login.
func waitForJWT() (string, error) {
	jwtChan := make(chan string)

	// Create a simple HTTP handler that listens for the JWT
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		jwt := r.URL.Query().Get("jwt")
		if jwt != "" {
			jwtChan <- jwt
			w.Write([]byte("Tapeless login successful! You can close this window."))
		} else {
			w.Write([]byte("Tapeless login error: JWT token missing!"))
		}
	})

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Waiting for JWT callback on http://localhost:%s/callback\n", env.LoginCallbackPort)
		http.ListenAndServe(":8080", nil)
	}()

	// Wait for the JWT to be received
	jwt := <-jwtChan
	return jwt, nil
}
