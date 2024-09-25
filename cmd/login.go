package cmd

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(loginCmd)
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
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
		fmt.Println("Waiting for JWT callback on http://localhost:8080/callback")
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
		loginURL := "http://localhost:5173/cli/login"
		fmt.Println("Opening browser to log in...")

		// Open browser for user to log in
		err := openBrowser(loginURL)
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
