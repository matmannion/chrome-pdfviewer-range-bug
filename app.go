package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting on http://localhost:8000/...")

	handleSameSiteCookies := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Don't allow caching
			w.Header().Add("Cache-Control", "no-cache")

			if r.URL.Path == "/" {
				// Set the header directly as go 1.10 doesn't support SameSite
				// Set two cookies, one without SameSite and one with
				w.Header().Add("Set-Cookie", "Auth1=auth1; SameSite=Lax")
				w.Header().Add("Set-Cookie", "Auth2=auth2")
			} else {
				// Read the two cookies. If we're missing either, serve a 404
				auth2, err2 := r.Cookie("Auth2")
				if err2 != nil || auth2.Value != "auth2" {
					w.Header().Add("Location", "/login")
					w.WriteHeader(http.StatusSeeOther)
					w.Write([]byte("303 - Missing auth2 cookie"))
					return
				}

				auth1, err1 := r.Cookie("Auth1")
				if err1 != nil || auth1.Value != "auth1" {
					w.Header().Add("Location", "/login")
					w.WriteHeader(http.StatusSeeOther)
					w.Write([]byte("303 - Missing auth1 cookie"))
					return
				}
			}

			h.ServeHTTP(w, r)
		}
	}

	http.Handle("/", logging(logger)(handleSameSiteCookies(http.FileServer(http.Dir("static")))))
	http.ListenAndServe(":8000", nil)
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.Header.Get("Range"), r.Header.Get("Cookie"))
			}()
			next.ServeHTTP(w, r)
		})
	}
}
