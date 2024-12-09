package main

import (
	"bytes"
	"context"
	"crawlero-app/db"
	"crawlero-app/routes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

var templates, _ = template.ParseFiles(
	"templates/layouts/console.html",
	"templates/pages/console/dashboard.html",
)

func init() {
	err := godotenv.Load()
	fmt.Println("Loading .env file")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

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

func uploadToS3(sess *session.Session, key string, body []byte) (string, error) {
	uploader := s3.New(sess)
	_, err := uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", os.Getenv("S3_BUCKET_URL"), key), nil
}

func consume() {
	reader := getKafkaReader(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_RESULT_TOPIC"),
		os.Getenv("GROUP_ID"),
	)
	defer reader.Close()

	sess, err := session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(os.Getenv("S3_REGION")),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
		Credentials:      credentials.NewStaticCredentials(os.Getenv("S3_KEY_ID"), os.Getenv("S3_SECRET_KEY"), ""),
	})

	if err != nil {
		log.Fatalf("failed to create a new session: %v", err)
	}

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		headers := m.Headers
		pool := db.GetDbPool()

		crawlerId := headers[0].Value
		runId := string(headers[1].Value)

		s3Key := fmt.Sprintf("crawl-results/%s/%s.json", crawlerId, runId)

		var data map[string]interface{}
		err = json.Unmarshal(m.Value, &data)
		if err != nil {
			log.Println("error unmarshalling message value")
			_, err = pool.Exec(
				context.Background(),
				"UPDATE runs SET status = 'failed' WHERE id = $1",
				runId,
			)
		}

		result, ok := data["result"]
		if !ok {
			log.Println("result key not found in message value")
			_, err = pool.Exec(
				context.Background(),
				"UPDATE runs SET status = 'failed' WHERE id = $1",
				runId,
			)
		}

		resultJson, _ := json.Marshal(result)
		s3Link, err := uploadToS3(sess, s3Key, resultJson)

		if err != nil {
			log.Println("error uploading to s3")
			_, err = pool.Exec(
				context.Background(),
				"UPDATE runs SET status = 'failed' WHERE id = $1",
				runId,
			)
		}

		_, err = pool.Exec(
			context.Background(),
			"UPDATE runs SET status = 'completed', result = $1 WHERE id = $2",
			s3Link,
			runId,
		)

		if err != nil {
			log.Fatalf("error updating run status: %v", err)
		}

		fmt.Println("updated run status")
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

	pool := db.GetDbPool()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
	}()

	http.ListenAndServe(":3000", r)

	defer pool.Close()
}
