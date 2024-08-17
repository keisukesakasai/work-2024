package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type Entity struct {
	Value     string
	CreatedAt time.Time
}

var (
	tracer = otel.Tracer("rolldice")
)

func main() {
	tracerProvider, err := InitTracer()
	if err != nil {
		log.Fatalf("Error setting up trace provider: %v", err)
	}
	defer func() { _ = tracerProvider.Shutdown(context.Background()) }()

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "Start Span")
	defer span.End()

	projectID := "datadog-sandbox"
	client, err := datastore.NewClient(
		ctx,
		projectID,
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	uuid := uuid.New().String()
	key := datastore.NameKey("keisukes-opentelemetry-work", uuid, nil)
	e := new(Entity)
	randomEmoji := getEmoji(ctx)
	e.Value = "Hello World " + randomEmoji
	loc, _ := time.LoadLocation("Asia/Tokyo")
	e.CreatedAt = time.Now().In(loc)

	if _, err := client.Put(ctx, key, e); err != nil {
		log.Fatalf("Failed to save entity: %v", err)
	}

	e = new(Entity)
	if err = client.Get(ctx, key, e); err != nil {
		log.Fatalf("Failed to get entity: %v", err)
	}

	fmt.Printf("Fetched entity: %v", e)

	client.Close()
}

func getEmoji(ctx context.Context) string {
	_, span := tracer.Start(ctx, "getEmoji")
	defer span.End()
	emojis := []string{"😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", "😊", "😇"}
	rand.NewSource(time.Now().UnixNano())
	randomEmoji := emojis[rand.Intn(len(emojis))]

	return randomEmoji
}
