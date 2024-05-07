package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	opentelemetry "app/internal"
)

var (
	tracer = otel.Tracer("otel-findy-demo")
)

func main() {
	// OpenTelemetry Traces
	tracerProvider, err := opentelemetry.InitTracer()
	if err != nil {
		log.Fatalf("Error setting up trace provider: %v", err)
	}
	defer func() { _ = tracerProvider.Shutdown(context.Background()) }()

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(mainHandler), "/")
	http.Handle("/", otelHandler)
	log.Fatalln(http.ListenAndServe(":18080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, span := tracer.Start(ctx, "main handler")
	defer span.End()

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	prosessing(ctx)
}

func prosessing(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "processing...")
	defer span.End()

	if rand.Float64() < 1.0/100.0 {
		funcAbnormal(ctx)
	} else {
		funcNormal(ctx)
	}
}

func funcNormal(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "funcNormal")
	defer span.End()

	request, err := http.NewRequestWithContext(ctx, "GET", "http://httpbin.org/get", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=============================")
	for name, values := range request.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}
	fmt.Println("=============================")

	time.Sleep(10 * time.Millisecond)
}

func funcAbnormal(ctx context.Context) {
	_, span := tracer.Start(ctx, "funcAbNormal(Oh...taking a lot of time...)")
	defer span.End()

	time.Sleep(3 * time.Second)
}
