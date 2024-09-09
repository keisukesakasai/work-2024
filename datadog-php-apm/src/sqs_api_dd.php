<?php

require '../vendor/autoload.php';

use Aws\Sqs\SqsClient;
use Aws\Exception\AwsException;
use Monolog\Logger;
use Monolog\Handler\StreamHandler;
use Monolog\Formatter\JsonFormatter;
use DDTrace\GlobalTracer;

/**
 * SQS クライアントを初期化する
 */
function initSqsClient($awsKey, $awsSecret) {
    return new SqsClient([
        'region'  => 'ap-northeast-1',
        'version' => 'latest',
        'credentials' => ['key' => $awsKey, 'secret' => $awsSecret],
    ]);
}

/**
 * ログのセットアップを行う
 */
function setupLogger() {
    $log = new Logger('php-sqs');
    $stream = new StreamHandler('/log/sqs-log.json', Logger::DEBUG);
    $stream->setFormatter(new JsonFormatter());
    $log->pushHandler($stream);
    return $log;
}

/**
 * メッセージを SQS に送信する
 */
function sendMessageToSqs($sqsClient, $queueUrl, $messageBody, $traceId, $parentId) {
    $messageAttributes = [
        '_datadog' => [
            'DataType' => 'String',
            'StringValue' => json_encode([
                'x-datadog-trace-id' => $traceId,
                'x-datadog-parent-id' => $parentId,
                // 'x-datadog-sampling-priority' => '1',
            ]),
        ],
    ];

    echo 'Message Attributes: ' . json_encode($messageAttributes, JSON_PRETTY_PRINT) . PHP_EOL;

    $result = $sqsClient->sendMessage([
        'QueueUrl' => $queueUrl,
        'MessageBody' => $messageBody,
        'MessageAttributes' => $messageAttributes,
    ]);

    return $result->get('MessageId');
}

// メイン処理
$awsKey = $_ENV['AWS_KEY'];
$awsSecret = $_ENV['AWS_SECRET'];
$queueUrl = 'https://sqs.ap-northeast-1.amazonaws.com/601427279990/sakasai-work-sqs';

$sqsClient = initSqsClient($awsKey, $awsSecret);
$log = setupLogger();

try {
    $tracer = GlobalTracer::get();
    $scope = $tracer->startActiveSpan('SQS send message');
    $span = $scope->getSpan();
    $traceContext = $span->getContext();
    
    $messageBody = json_encode(['message' => 'Hello from PHP!', 'timestamp' => time()]);
    echo 'Message Body: ' . $messageBody . PHP_EOL;

    $messageId = sendMessageToSqs($sqsClient, $queueUrl, $messageBody, $traceContext->getTraceId(), $traceContext->getSpanId());

    $span->finish();

    $log->info("Message sent successfully", ['MessageId' => $messageId]);
    echo 'Message sent successfully. MessageId: ' . $messageId . PHP_EOL;

} catch (AwsException $e) {
    $log->error($e->getAwsErrorMessage());
    echo 'AWS Error: ' . $e->getAwsErrorMessage() . PHP_EOL;
} catch (Exception $e) {
    $log->error($e->getMessage());
    echo 'Error: ' . $e->getMessage() . PHP_EOL;
}
