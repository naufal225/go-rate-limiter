package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/naufal225/go-inventory-api/helpers"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /ping", helpers.RateLimitMiddleware(http.HandlerFunc(Ping)))

	fmt.Println("Server start serving and listen on :8080")

	go func ()  {
		for true {
			time.Sleep(5 * time.Second)	
			helpers.Mu.Lock()
			for k, v := range helpers.IPCount {
				if(time.Since(v.WindowStart).Seconds() > 60) {
					delete(helpers.IPCount, k)
				}
			}
			helpers.Mu.Unlock()
		}
	}()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{
		"message": "pong",
	}

	jsonInBytes, err := json.Marshal(&message)
	if err != nil {
		fmt.Println("Gagal konversi ke json")
		return
	}

	w.Write(jsonInBytes)
}
