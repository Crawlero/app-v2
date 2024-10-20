package routes

import (
	"context"
	"crawlero-app/db"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
)

var consoleLayout, _ = template.ParseFiles(
	"templates/layouts/console.html",
)

var listCrawlerTmpl, _ = template.Must(consoleLayout.Clone()).ParseFiles(
	"templates/pages/console/crawler/index.html",
)

var detailCrawlerTmpl, _ = template.Must(consoleLayout.Clone()).ParseFiles(
	"templates/pages/console/crawler/detail.html",
)

type Crawler struct {
	ID     string
	Name   string
	Url    string
	Status string
}

func CrawerRoutes() chi.Router {
	router := chi.NewRouter()
	router.Get("/", listCrawler)
	router.Post("/", createCrawler)
	router.Get("/{crawlerID}", getOneCrawler)

	// router.Route("/{articleID}", func(r chi.Router) {
	// 	r.Use(ar.ArticleCtx)            // Load the *Article on the request context
	// 	r.Get("/", ar.GetArticle)       // GET /articles/123
	// 	r.Put("/", ar.UpdateArticle)    // PUT /articles/123
	// 	r.Delete("/", ar.DeleteArticle) // DELETE /articles/123
	// })
	return router
}

func listCrawler(w http.ResponseWriter, r *http.Request) {
	pool := db.GetDbPool()
	var crawlers []Crawler

	rows, err := pool.Query(context.Background(), "SELECT id, name, url, status FROM crawlers")

	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("Result:", rows)

	for rows.Next() {
		var id string
		var name string
		var url string
		var status string

		err = rows.Scan(&id, &name, &url, &status)
		if err != nil {
			fmt.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		crawlers = append(crawlers, Crawler{
			ID:     id,
			Name:   name,
			Url:    url,
			Status: status,
		})
	}

	data := map[string]interface{}{
		"Crawlers": crawlers,
	}

	if err := listCrawlerTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getOneCrawler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "Crawler",
		"Message": "Crawler",
	}

	if err := detailCrawlerTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createCrawler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Crawler")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := r.FormValue("name")
	desc := r.FormValue("description")

	fmt.Println("Name:", name)
	fmt.Println("Desc:", desc)

	http.Redirect(w, r, "/crawler", http.StatusSeeOther)
}
