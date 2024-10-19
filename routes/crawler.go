package routes

import (
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
	Url    int
	Status string
}

func CrawerRoutes() chi.Router {
	router := chi.NewRouter()
	router.Get("/", listCrawler)
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
	crawlers := []Crawler{
		{Name: "John Brown", Url: 45, Status: "active", ID: "1"},
		{Name: "Jim Green", Url: 27, Status: "inactive", ID: "2"},
		{Name: "Joe Black", Url: 31, Status: "active", ID: "3"},
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
