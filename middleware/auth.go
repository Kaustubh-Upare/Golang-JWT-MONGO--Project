package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func TestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Println("Hello I am Middleware Just testing", start)
		admino := r.URL.Query().Get("admin")

		if admino == "" {
			http.Error(w, "Invalid Access", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))

	})
}
