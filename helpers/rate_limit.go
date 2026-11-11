package helpers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimitEntry struct {
	Count int
	WindowStart time.Time
}

var IPCount map[string]*RateLimitEntry = make(map[string]*RateLimitEntry)
var Limit int = 5
var Mu sync.Mutex

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println(err)
		}

		var kenaLimit = false

		Mu.Lock()

		rateLimitEntry, ok := IPCount[host]
		if !ok || time.Since(rateLimitEntry.WindowStart).Seconds() > 5 {
			IPCount[host] = &RateLimitEntry{
				Count: 1,
				WindowStart: time.Now(),
			}
		} else {
			if rateLimitEntry.Count >= Limit {
				kenaLimit = true
			}

			IPCount[host].Count = rateLimitEntry.Count + 1
		}

		Mu.Unlock()

		if kenaLimit {
			fmt.Println("Kena limit")

			var message map[string]string = map[string]string{
				"message": "kamu telah mencapai limit, silahkan coba lagi nanti",
			}

			jsonInBytes, err := json.Marshal(&message)
			if err != nil {
				fmt.Println("Gagal konversi ke json")
				return
			}

			w.WriteHeader(http.StatusTooManyRequests)
			w.Write(jsonInBytes)
			return
		}
		next.ServeHTTP(w, r)
	})
}
