package main

import (
    "encoding/json"
    "net/http"
    "github.com/sirupsen/logrus"
)

func main() {
    log := logrus.New()
    // Set the formatter to include the time and date.
    log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })

    http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        response := map[string]string{"message": "This is the Product Service"}
        json.NewEncoder(w).Encode(response)
    })

    log.Info("Product service starting on port 8081")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal("Failed to start server: ", err)
    }
}
