package main

import (
	"context"
	"crawlero-app/routes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/segmentio/kafka-go"
)

var templates, _ = template.ParseFiles(
	"templates/layouts/console.html",
	"templates/pages/console/dashboard.html",
)

const (
	KAFKA_RESULT_TOPIC = "crawl_results"
	KAFKA_ERROR_TOPIC  = "crawl_errors"

	KAFKA_BROKER = "127.0.0.1:9093"

	GROUP_ID = "crawlero-app"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func consume() {
	reader := getKafkaReader(KAFKA_BROKER, KAFKA_RESULT_TOPIC, GROUP_ID)
	defer reader.Close()

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

        headers := m.Headers
        for _, header := range headers {
            fmt.Printf("Header: %s = %s\n", header.Key, header.Value)
        }
	}
}

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

    go consume()

	http.ListenAndServe(":3000", r)
}
