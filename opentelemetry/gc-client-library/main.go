package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type Entity struct {
	Value     string
	CreatedAt time.Time
}

var (
	tracer     = otel.Tracer("store-emoji")
	projectID  = "datadog-sandbox"
	databaseID = "sakasaik-datastore"
)

func main() {
	tracerProvider, err := InitTracer()
	if err != nil {
		log.Fatalf("Error setting up trace provider: %v", err)
	}
	defer func() { _ = tracerProvider.Shutdown(context.Background()) }()

	http.Handle("/store-emoji", otelhttp.NewHandler(http.HandlerFunc(storeEmojiHandler), "storeEmojiHandler"))

	server := &http.Server{
		Addr: ":8080",
	}

	// サーバーを並行処理で開始
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// シグナルハンドリング
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// シグナルが来るまで待機
	<-stop

	// サーバーのシャットダウン処理
	log.Println("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	wg.Wait()
	log.Println("Server exited")
}

func storeEmojiHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "storeEmojiHandler")
	log.Printf("Trace ID: %s", span.SpanContext().TraceID())
	defer span.End()

	client, err := datastore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create datastore client: %v", err), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	uuid := uuid.New().String()
	key := datastore.NameKey("keisukes-opentelemetry-work", uuid, nil)

	e := &Entity{
		Value:     "Hello World " + getEmoji(ctx),
		CreatedAt: time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)),
	}

	if _, err := client.Put(ctx, key, e); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save entity: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Stored entity with key %s and value: %s\n", uuid, e.Value)
}

func getEmoji(ctx context.Context) string {
	_, span := tracer.Start(ctx, "getRundomEmoji")
	defer span.End()
	emojis := []string{"😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", "😊", "😇"}
	rand.NewSource(time.Now().UnixNano())
	randomEmoji := emojis[rand.Intn(len(emojis))]

	return randomEmoji
}
