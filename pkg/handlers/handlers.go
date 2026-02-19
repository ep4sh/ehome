package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"ehome/pkg/app"
)

func Index(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("templates/base.html",
			"templates/head.html",
			"templates/topmenu.html",
			"templates/footer.html",
			"templates/index.html")
		if err != nil {
			log.Printf("error parsing template files: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(w, "base.html")
		if err != nil {
			log.Printf("error executing template files: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Static(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("static"))
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	}
}

func Healthz(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func Readyz(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// LoggingMiddleware wraps every request with logging data.
func LoggingMiddleware(app *app.Application, next http.Handler) http.Handler {
	start := time.Now()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			app.Sugar.Infow("Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
}

// SecurityHeadersMiddleware wraps every request with security headers.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Strict-Transport-Security", "max-age=2678400; includeSubDomains")
			w.Header().Set("Feature-Policy", "self")
			w.Header().Set("X-Frame-Options", "DENY")
			next.ServeHTTP(w, r)
		})
}
