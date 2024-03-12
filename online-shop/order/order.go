package main

import (
    "context"
	"crypto/rand"
    "fmt"
    "net/http"
	"os"
	"os/exec"
    "runtime"
	"encoding/json"
	"encoding/base64"

    "github.com/sirupsen/logrus"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var (
    //Callback URL
    redirectURL = "http://localhost:8082/callback"
)

// generateState generates a random string
func generateState() string {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

// openBrowser tries to open the browser with a given URL.
func openBrowser(url string) {
    var err error
	
    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        logrus.Printf("Please open your web browser and visit: %s\n", url)
    }
    if err != nil {
        logrus.Printf("Failed to open browser: %s\n", err)
    }
}

func main() {
    log := logrus.New()
    log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })

    ctx := context.Background()
    conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
        ClientSecret: os.Getenv("CLIENT_SECRET"),
        RedirectURL:  redirectURL,
        Scopes:       []string{"openid", "email"},
        Endpoint:     google.Endpoint,
    }

	state := generateState()

    // OIDC callback handler
    http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Query().Get("state") != state {
            http.Error(w, "state did not match", http.StatusBadRequest)
            return
        }

        code := r.URL.Query().Get("code")
        token, err := conf.Exchange(ctx, code)
        if err != nil {
            http.Error(w, "failed to exchange token: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Token retrieved
        log.Infof(w, "Authentication successful! Token: %v", token)
		
		// Print the ID token (JWT) to the console
		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "Failed to get id_token", http.StatusInternalServerError)
			return
		}
		log.Infof("Authentication successful! ID Token: %s", idToken)
    })

    // Order service endpoint
    http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        response := map[string]string{"message": "This is the Order Service"}
        json.NewEncoder(w).Encode(response)
    })

    // Start the authentication flow by redirecting the user to the Google login page
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        authURL := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
        log.Infof("You will be redirected to the Google login page. If not, open the following URL in your browser:", authURL)
        openBrowser(authURL) // Call to corrected openBrowser function
        http.Redirect(w, r, authURL, http.StatusFound)
    })

    log.Infof("Order service starting on port 8082")
    if err := http.ListenAndServe(":8082", nil); err != nil {
        fmt.Printf("Failed to start server: %s\n", err)
        os.Exit(1)
    }
}