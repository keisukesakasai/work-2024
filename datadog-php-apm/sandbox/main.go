package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/sirupsen/logrus"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// SQSクライアントを初期化する
func initSqsClient(awsKey, awsSecret, region string) (*sqs.Client, error) {
	// 認証情報プロバイダを設定
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(awsKey, awsSecret, ""))

	// AWS設定をロード
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}

// メッセージを SQS に送信する
func sendMessageToSqs(sqsClient *sqs.Client, queueUrl string, messageBody string) (string, error) {
	ctx := context.Background()

	// 現在のスパンをコンテキストから取得
	span, _ := ddtracer.SpanFromContext(ctx)
	traceID := span.Context().TraceID()
	spanID := span.Context().SpanID()

	messageAttributes := map[string]types.MessageAttributeValue{
		"_datadog": {
			DataType:    aws.String("String"),
			StringValue: aws.String(fmt.Sprintf(`{"trace_id":"%d","span_id":"%d"}`, traceID, spanID)),
		},
	}

	result, err := sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:          &queueUrl,
		MessageBody:       &messageBody,
		MessageAttributes: messageAttributes,
	})
	if err != nil {
		return "", err
	}

	return *result.MessageId, nil
}

func main() {
	awsKey := os.Getenv("AWS_KEY")
	awsSecret := os.Getenv("AWS_SECRET")
	queueUrl := "https://sqs.ap-northeast-1.amazonaws.com/601427279990/sakasai-work-sqs"
	region := "ap-northeast-1"

	sqsClient, err := initSqsClient(awsKey, awsSecret, region)
	if err != nil {
		fmt.Printf("Failed to initialize SQS client: %v", err)
	}

	ddtracer.Start(ddtracer.WithServiceName("go-sqs-app"))
	defer ddtracer.Stop()

	span, _ := ddtracer.StartSpanFromContext(context.Background(), "SQS send message")
	defer span.Finish()

	messageBody, err := json.Marshal(map[string]interface{}{
		"message":   "Hello from Go!",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		fmt.Printf("Failed to marshal message body: %v", err)
	}

	fmt.Println("Message Body: ", string(messageBody))

	messageId, _ := sendMessageToSqs(sqsClient, queueUrl, string(messageBody))
	fmt.Println("Message sent successfully", logrus.Fields{"MessageId": messageId})
	fmt.Println("Message sent successfully. MessageId: ", messageId)
}
