package routes

import (
	"context"
	"crawlero-app/db"
	"crawlero-app/job"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

var crawlerRunTmpl, _ = template.Must(consoleLayout.Clone()).ParseFiles(
	"templates/parts/crawler-detail-header.html",
	"templates/parts/run_row.html",
	"templates/pages/console/crawler/run.html",
)

var crawlerRunRowTmpl, _ = template.ParseFiles(
	"templates/parts/run_row.html",
)

var crawlerRunDetailTmpl, _ = template.Must(consoleLayout.Clone()).ParseFiles(
    "templates/parts/crawler-detail-header.html",
	"templates/pages/console/crawler/run-detail.html.tmpl",
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
	router.Put("/{crawlerID}/schema", updateCrawlerSchema)
	router.Get("/{crawlerID}/schedule", getCrawlerSchedule)
	router.Get("/{crawlerID}/run", getListRun)
	router.Post("/{crawlerID}/run", createRun)
	router.Get("/{crawlerID}/run/{runID}", runDetail)

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

	var crawlerSchema sql.NullString

	err := pool.QueryRow(context.Background(), "SELECT schema FROM crawlers WHERE id = $1", crawlerID).Scan(
		&crawlerSchema,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := crawlerSchemaTmpl.Execute(
		w,
		map[string]interface{}{
			"ID":     crawlerID,
			"Schema": crawlerSchema,
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

func updateCrawlerSchema(w http.ResponseWriter, r *http.Request) {
	crawlerID := chi.URLParam(r, "crawlerID")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("Form Data:", r.Form)
	fmt.Println(r.FormValue("schema"))
	fmt.Println(r.FormValue("n"))
	fmt.Println("Okok")

	pool := db.GetDbPool()
	_, err := pool.Exec(
		context.Background(),
		"UPDATE crawlers SET schema = $1 WHERE id = $2",
		r.FormValue("schema"),
		crawlerID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Updated"))
}

func getListRun(w http.ResponseWriter, r *http.Request) {
	crawlerID := chi.URLParam(r, "crawlerID")
	pool := db.GetDbPool()

	rows, err := pool.Query(
		context.Background(),
		"SELECT id, status, started_at, updated_at FROM runs WHERE crawler_id = $1 ORDER BY updated_at DESC",
		crawlerID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var runs []map[string]interface{}

	for rows.Next() {
		var id string
		var status string
		var createdAt pgtype.Timestamp
		var updatedAt pgtype.Timestamp

		err = rows.Scan(&id, &status, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		runs = append(runs, map[string]interface{}{
			"ID":        id,
			"Status":    status,
			"CreatedAt": createdAt.Time.Format("2006-01-02 15:04"),
			"UpdatedAt": updatedAt.Time.Format("2006-01-02 15:04"),
			"CrawlerID": crawlerID,
		})
	}

	if err := crawlerRunTmpl.Execute(
		w,
		map[string]interface{}{
			"ID":   crawlerID,
			"Runs": runs,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createRun(w http.ResponseWriter, r *http.Request) {
	crawlerID := chi.URLParam(r, "crawlerID")
	pool := db.GetDbPool()

	var schemaStr sql.NullString
	var url string

	err := pool.QueryRow(context.Background(), "SELECT schema, url FROM crawlers WHERE id = $1", crawlerID).Scan(
		&schemaStr,
		&url,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !schemaStr.Valid {
		http.Error(w, "Schema is not defined", http.StatusBadRequest)
		return
	}
	if url == "" {
		http.Error(w, "URL is not defined", http.StatusBadRequest)
		return
	}

	schema := job.Schema{}
	err = json.Unmarshal([]byte(schemaStr.String), &schema)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	runInput := job.RunInput{
		Url:    url,
		Schema: schema,
	}

	var runId string
	err = pool.QueryRow(
		context.Background(),
		"INSERT INTO runs (crawler_id, status, input, started_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING id",
		crawlerID,
		"running",
		runInput,
	).Scan(&runId)

	err = job.CreateCrawJob(crawlerID, runId, runInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/crawler/%s/run", crawlerID))
	w.WriteHeader(http.StatusCreated)
}

func runDetail(w http.ResponseWriter, r *http.Request) {
	runID := chi.URLParam(r, "runID")
	pool := db.GetDbPool()

	runDetail := struct {
		CrawlerId string
		RunId     string
		Status    string
		Input     string
		Result    string
		StartedAt pgtype.Timestamp
		UpdatedAt pgtype.Timestamp
	}{}

	err := pool.QueryRow(
		context.Background(),
		"SELECT id, crawler_id, status, input, result, started_at, updated_at FROM runs WHERE id = $1",
		runID,
	).Scan(
		&runDetail.RunId,
		&runDetail.CrawlerId,
		&runDetail.Status,
		&runDetail.Input,
		&runDetail.Result,
		&runDetail.StartedAt,
		&runDetail.UpdatedAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := crawlerRunDetailTmpl.Execute(
		w,
		map[string]interface{}{
			"ID": runDetail.CrawlerId,
			"Run": map[string]interface{}{
				"ID":        runDetail.RunId,
				"Status":    runDetail.Status,
				"Input":     runDetail.Input,
				"Result":    runDetail.Result,
				"StartedAt": runDetail.StartedAt.Time.Format("2006-01-02 15:04"),
				"UpdatedAt": runDetail.UpdatedAt.Time.Format("2006-01-02 15:04"),
			},
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
