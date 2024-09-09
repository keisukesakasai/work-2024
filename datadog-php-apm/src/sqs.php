<?php

require '../vendor/autoload.php';

use Aws\Sqs\SqsClient;
use Aws\Exception\AwsException;
use Monolog\Logger;
use Monolog\Handler\StreamHandler;
use Monolog\Formatter\JsonFormatter;
use DDTrace\GlobalTracer;

$awsKey = $_ENV['AWS_KEY'];
$awsSecret = $_ENV['AWS_SECRET'];

// AWS SQS クライアントを初期化
$sqsClient = new SqsClient([
    'region'  => 'ap-northeast-1',
    'version' => 'latest',
    'credentials' => [
        'key'    => $awsKey,
        'secret' => $awsSecret,
    ],
]);

// SQS キュー URL
$queueUrl = 'https://sqs.ap-northeast-1.amazonaws.com/601427279990/sakasai-work-sqs'; // 適切なキューURLに置き換える

// ログのセットアップ
$log = new Logger('php-sqs');
$formatter = new JsonFormatter();
$stream = new StreamHandler('/log/sqs-log.json', Logger::DEBUG);
$stream->setFormatter($formatter);
$log->pushHandler($stream);

try {
    // Datadog トレーサーを取得
    $tracer = GlobalTracer::get();
    
    // スパンを作成 (SQS メッセージ送信用)
    $scope = $tracer->startActiveSpan('SQS send message');
    $span = $scope->getSpan();
    
    // トレースコンテキストを抽出
    $traceContext = $span->getContext();
    $traceId = $traceContext->getTraceId();
    $parentId = $traceContext->getSpanId();

    // メッセージ本体
    $messageBody = json_encode([
        'message' => 'Hello from PHP!',
        'timestamp' => time(),
    ]);

    echo 'Message Body: ' . $messageBody . PHP_EOL;

    // メッセージを送信
    /*
    $result = $sqsClient->sendMessage([
        'QueueUrl'    => $queueUrl,
        'MessageBody' => $messageBody,
        'MessageAttributes' => [
            'x-datadog-trace-id' => [
                'DataType' => 'String',
                'StringValue' => (string) $traceId,
            ],
            'x-datadog-parent-id' => [
                'DataType' => 'String',
                'StringValue' => (string) $parentId,
            ],
        ],
    ]);
    */

    // メッセージを単純に送信
    $result = $sqsClient->sendMessage([
        'QueueUrl'    => $queueUrl,
        'MessageBody' => $messageBody,
    ]);

    // スパンを終了
    $span->finish();

    $log->info("Message sent successfully", ['MessageId' => $result->get('MessageId')]);
    echo 'Message sent successfully. MessageId: ' . $result->get('MessageId') . PHP_EOL;

} catch (AwsException $e) {
    $errorMessage = 'AWS Error: ' . $e->getAwsErrorMessage() . ' (Error Type: ' . $e->getAwsErrorType() . ')';
    $log->error($errorMessage);
    echo 'AWS Error Code: ' . $e->getAwsErrorCode() . PHP_EOL;
    echo 'AWS Error Message: ' . $e->getAwsErrorMessage() . PHP_EOL;
    echo 'AWS Error Type: ' . $e->getAwsErrorType() . PHP_EOL;
} catch (Exception $e) {
    $errorMessage = 'General Error: ' . $e->getMessage();
    $log->error($errorMessage);
    echo $errorMessage . PHP_EOL;
}
