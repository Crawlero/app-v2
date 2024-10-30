package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

const (
	KAFKA_BROKER      = "127.0.0.1:9093"
	KAFKA_INPUT_TOPIC = "crawl_inputs"
)

type Field struct {
	Name      string  `json:"name"`
	Selector  string  `json:"selector,omitempty"`
	Type      string  `json:"type"`
	Attribute string  `json:"attribute,omitempty"`
	Fields    []Field `json:"fields,omitempty"`
}

type Schema struct {
	Name         string  `json:"name"`
	BaseSelector string  `json:"baseSelector"`
	Fields       []Field `json:"fields"`
}

type RunInput struct {
	Url    string `json:"url"`
	Schema Schema `json:"schema"`
}

func CreateCrawJob(crawlerId string, runId string, input RunInput) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{KAFKA_BROKER},
		Topic:   KAFKA_INPUT_TOPIC,
	})
	defer func() {
		if err := writer.Close(); err != nil {
			log.Fatalf("Failed to close writer: %v", err)
		}
	}()

	inputJson, _ := json.Marshal(input)
	msg := kafka.Message{
		Key:   []byte(runId),
		Value: inputJson,
		Headers: []kafka.Header{
			{Key: "crawlerId", Value: []byte(crawlerId)},
			{Key: "runId", Value: []byte(runId)},
		},
	}

	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println("Error while writing message to Kafka.")
		fmt.Println(err)
		return err
	}

	fmt.Println("Message written to Kafka successfully.")
	return nil
}
