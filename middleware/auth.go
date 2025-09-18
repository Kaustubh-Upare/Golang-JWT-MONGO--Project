package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kaustubh-upare/jwtWithMongo/utils"
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, "You Need to login or Signin First", http.StatusForbidden)
				return
			}
			http.Error(w, "You Need to login or Signin First", http.StatusInternalServerError)
			return
		}
		if cookie.Value == "" {
			http.Error(w, "Dont Play with Cookies", http.StatusInternalServerError)
			return
		}

		err = utils.ValidateCookie(cookie.Value)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid or Expired Token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
