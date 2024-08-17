package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	_ "time/tzdata"

	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/exp/rand"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gorilla/mux"
)

var tracer = otel.Tracer("datadog-otel-lambda")

func main() {
	// opentelemetry
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// http server
	mux := mux.NewRouter()
	mux.Use(otelmux.Middleware("datadog-otel-lambda"))
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// load aws sdk
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion("ap-northeast-1"),
		)
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		// instrument all aws clients
		otelaws.AppendMiddlewares(&cfg.APIOptions)

		// for lambda
		svc := lambda.NewFromConfig(cfg)
		functionName1 := os.Getenv("FUNCTION_NAME_1")
		functionName2 := os.Getenv("FUNCTION_NAME_2")

		// create data
		lengthStr := r.URL.Query().Get("length")
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			http.Error(w, "Invalid length parameter", http.StatusBadRequest)
		}
		data := createData(ctx, length)
		fmt.Println("bitstring: ", data)

		// invoke lambda with datadog extension layer
		response1, err := invokeLambdaFunction(ctx, svc, functionName1, data)
		if err != nil {
			log.Printf("Error invoking function %s: %v", functionName1, err)
		} else {
			fmt.Printf("Response from %s: %v\n", functionName1, response1)
		}

		// invoke lambda with otel custom extension layer
		response2, err := invokeLambdaFunction(ctx, svc, functionName2, data)
		if err != nil {
			log.Printf("Error invoking function %s: %v", functionName2, err)
		} else {
			fmt.Printf("Response from %s: %v\n", functionName2, response2)
		}

		response := map[string]interface{}{
			"response1": response1,
			"response2": response2,
		}

		// レスポンスをエンコードして返却
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	serverPort := os.Getenv("SERVER_PORT")
	fmt.Println("Server Port: ", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

func invokeLambdaFunction(ctx context.Context, svc *lambda.Client, functionName string, data string) (map[string]interface{}, error) {
	ctx, span := tracer.Start(ctx, "invokeLambdaFunction: "+functionName)
	defer span.End()

	// encode data
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	carrierBytes, _ := json.Marshal(carrier)

	payload := map[string]string{
		"bitstring": data,
		"carrier":   string(carrierBytes),
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// invoke lambda
	result, err := svc.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &functionName,
		Payload:        jsonPayload,
		InvocationType: types.InvocationTypeRequestResponse,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke function %s: %w", functionName, err)
	}

	// decode data
	var responsePayload map[string]interface{}
	err = json.Unmarshal(result.Payload, &responsePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response payload: %w", err)
	}

	return responsePayload, nil
}

func createData(ctx context.Context, length int) string {
	_, span := tracer.Start(ctx, "createData")
	defer span.End()
	rand.Seed(uint64(time.Now().UnixNano()))
	data := make([]byte, length)
	for i := range data {
		data[i] = byte('0' + rand.Intn(2)) // '0' or '1'
	}
	return string(data)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
