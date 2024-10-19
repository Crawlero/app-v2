package main

import (
	"crawlero-app/routes"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var templates, _ = template.ParseFiles(
	"templates/layouts/console.html",
	"templates/pages/console/dashboard.html",
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Dashboard
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"Title":   "Hello World",
			"Message": "Hello, World!",
		}
		if err := templates.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

    // Mounting routes
	r.Mount("/crawler", routes.CrawerRoutes())

	http.ListenAndServe(":3000", r)
}
