package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	log := logrus.New()
    // Set the formatter to include the time and date.
    log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })

	// Handler for the frontend service
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "This is the Frontend Service"}
		json.NewEncoder(w).Encode(response)
	})

	// Start ticker to check the status of cart and product services every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			checkService("http://online-shop-cart.online-shop.svc.cluster.local:8083/cart")
			checkService("http://online-shop-product.online-shop.svc.cluster.local:8081/product")
		}
	}()

	// Start the server
	log.Info("Frontend service starting on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

// checkService makes an HTTP GET request to the specified URL and logs the result using logrus
func checkService(url string, log *logrus.Logger) {
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error accessing %s: %s", url, err)
		return
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 (OK)
	if resp.StatusCode == http.StatusOK {
		log.Infof("Service %s is OK", url)
	} else {
		log.Warnf("Service %s returned HTTP error: %s", url, resp.Status)
	}
}
