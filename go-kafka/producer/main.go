package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	_ "time/tzdata"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func main() {
	// kafka 設定
	brokerList := []string{"kafka:9092"}
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	// Http server
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		// リクエストのボディを取得します
		_, err := io.ReadAll(c.Request.Body)
		if err != nil {
			http.Error(c.Writer, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// kafka producer
		producer, err := newAccessLogProducer(brokerList)
		if err != nil {
			log.Fatal(err)
		}

		// 送信するメッセージを作成します
		topic := "topic-otel"
		msg := sarama.ProducerMessage{
			Topic: topic,
		}

		// メッセージを送信します
		producer.Input() <- &msg
		successMsg := <-producer.Successes()
		log.Printf("Message sent topic: %s successfully! Partition: %d, Offset: %d", topic, successMsg.Partition, successMsg.Offset)

		err = producer.Close()
		if err != nil {
			log.Fatalln("Failed to close producer:", err)
		}
	})

	r.Run(":8080")
}

func newAccessLogProducer(brokerList []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("starting Sarama producer: %w", err)
	}

	go func() {
		for {
			select {
			case err := <-producer.Errors():
				log.Printf("Failed to write message: %v", err)
			case success := <-producer.Successes():
				log.Printf("Message sent successfully: %v", success)
			}
		}
	}()

	// Send a test message
	message := &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("test message"),
	}

	producer.Input() <- message

	return producer, nil
}
