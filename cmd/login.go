package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/env"
	"tapeless.app/tapeless-cli/util"
)

func init() {
	RootCmd.AddCommand(loginCmd)
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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Tapeless",
	Run: func(cmd *cobra.Command, args []string) {
		loginURL := env.WebURL + "/cli/login"
		fmt.Println("Opening browser to log in...")

		// Open browser for user to log in
		err := util.OpenBrowser(loginURL)
		if err != nil {
			fmt.Println("Error opening browser:", err)
			return
		}

		// Optional: Wait for JWT (via a callback server or polling mechanism)
		jwt, err := waitForJWT()
		if err != nil {
			fmt.Println("Error fetching JWT:", err)
			return
		}

		// Store the JWT in your config
		viper.Set("token", jwt)
		viper.WriteConfig()
		fmt.Println("Login successful:", jwt)
	},
}
