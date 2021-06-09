package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/tsawler/bookings-app/internal/helpers"
)

// NoSurf is the csrf protection middleware
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves session data for current request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

//check the users is login or not, if not login ,just throw err
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Please login first")
			http.Redirect(w,r,"/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w,r)
	})
}