package routes

import (
	"context"
	"crawlero-app/db"
	"database/sql"
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

var crawlerInfoTmpl, _ = template.Must(consoleLayout.Clone()).ParseFiles(
	"templates/parts/crawler-detail-header.html",
	"templates/pages/console/crawler/detail.html",
)

var crawlerSchemaTmpl, _ = template.Must(consoleLayout.Clone()).ParseFiles(
	"templates/parts/crawler-detail-header.html",
	"templates/pages/console/crawler/schema.html",
)

var crawlerScheduleTmpl, _ = template.Must(consoleLayout.Clone()).ParseFiles(
	"templates/parts/crawler-detail-header.html",
	"templates/pages/console/crawler/schedule.html",
)

type Crawler struct {
	ID          string
	Name        string
	Url         string
	Status      string
	Description sql.NullString
}

func CrawerRoutes() chi.Router {
	router := chi.NewRouter()
	router.Get("/", listCrawler)
	router.Post("/", createCrawler)
	router.Get("/{crawlerID}", getCrawlerInfo)
	router.Put("/{crawlerID}", updateCrawlerInfo)
	router.Get("/{crawlerID}/schema", getCrawlerSchema)
	router.Get("/{crawlerID}/schedule", getCrawlerSchedule)

	return router
}

func listCrawler(w http.ResponseWriter, r *http.Request) {
	pool := db.GetDbPool()
	var crawlers []Crawler

	rows, err := pool.Query(context.Background(), "SELECT id, name, url, status FROM crawlers")

	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var id string
		var name string
		var url string
		var status string

		err = rows.Scan(&id, &name, &url, &status)
		if err != nil {
			fmt.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

func getCrawlerInfo(w http.ResponseWriter, r *http.Request) {
	crawlerID := chi.URLParam(r, "crawlerID")
	pool := db.GetDbPool()
	crawler := Crawler{}

	err := pool.QueryRow(context.Background(), "SELECT id, name, url, status, description FROM crawlers WHERE id = $1", crawlerID).Scan(
		&crawler.ID,
		&crawler.Name,
		&crawler.Url,
		&crawler.Status,
		&crawler.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := crawlerInfoTmpl.Execute(w, crawler); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getCrawlerSchema(w http.ResponseWriter, r *http.Request) {
	crawlerID := chi.URLParam(r, "crawlerID")
	pool := db.GetDbPool()
	crawler := Crawler{}

	err := pool.QueryRow(context.Background(), "SELECT id, name, url, status, description FROM crawlers WHERE id = $1", crawlerID).Scan(
		&crawler.ID,
		&crawler.Name,
		&crawler.Url,
		&crawler.Status,
		&crawler.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fields := []struct{ Name, Selector, Type string }{
		{Name: "Title", Selector: ".title", Type: "text"},
		{Name: "Link", Selector: ".link", Type: "url"},
	}

	if err := crawlerSchemaTmpl.Execute(
		w,
		map[string]interface{}{
			"ID":     crawlerID,
			"Fields": fields,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getCrawlerSchedule(w http.ResponseWriter, r *http.Request) {
	crawlerID := chi.URLParam(r, "crawlerID")
	pool := db.GetDbPool()
	crawler := Crawler{}

	err := pool.QueryRow(context.Background(), "SELECT id, name, url, status, description FROM crawlers WHERE id = $1", crawlerID).Scan(
		&crawler.ID,
		&crawler.Name,
		&crawler.Url,
		&crawler.Status,
		&crawler.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := crawlerScheduleTmpl.Execute(
		w,
		crawler,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createCrawler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := r.FormValue("name")
	desc := r.FormValue("description")

	pool := db.GetDbPool()
	_, err := pool.Exec(context.Background(), "INSERT INTO crawlers (name, url, status) VALUES ($1, $2, $3)", name, desc, "active")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/crawler", http.StatusSeeOther)
}

func updateCrawlerInfo(w http.ResponseWriter, r *http.Request) {
	crawlerID := chi.URLParam(r, "crawlerID")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := r.FormValue("name")
	desc := r.FormValue("description")
	url := r.FormValue("url")

	pool := db.GetDbPool()
	_, err := pool.Exec(
		context.Background(),
		"UPDATE crawlers SET name = $1, url = $2, description = $3 WHERE id = $4",
		name,
		url,
		desc,
		crawlerID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Updated"))
}
