<?php

require '../vendor/autoload.php';

use Aws\Sqs\SqsClient;
use Aws\Exception\AwsException;
use Monolog\Logger;
use Monolog\Handler\StreamHandler;
use Monolog\Formatter\JsonFormatter;
use DDTrace\GlobalTracer;

/**
 * 大きな数値を 16 進数に変換するための関数 (bcmath 使用)
 */
function bcdechex($dec) {
    $hex = '';
    while (bccomp($dec, '0') > 0) {
        $remainder = bcmod($dec, '16');
        $hex = dechex($remainder) . $hex;
        $dec = bcdiv($dec, '16', 0);
    }
    return $hex;
}

/**
 * W3C Trace Context 形式で `traceparent` を生成する関数
 */
function create_traceparent($traceId, $parentId, $sampled = '01') {
    // bc で trace_id を 16進数に変換し、128 ビットに拡張
    $traceId_128bit = str_pad(bcdechex($traceId), 32, '0', STR_PAD_LEFT);

    // bc で span_id を 16進数に変換し、64 ビットに変換
    $parentId_64bit = str_pad(bcdechex($parentId), 16, '0', STR_PAD_LEFT);

    // W3C Trace Context の version は "00"
    $version = '00';

    // W3C Trace Context 形式の traceparent を生成
    $traceparent = sprintf('%s-%s-%s-%s', $version, $traceId_128bit, $parentId_64bit, $sampled);

    return $traceparent;
}

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

    // traceparent を作成
    $traceparent = create_traceparent($traceId, $parentId);
    echo "Generated traceparent: $traceparent" . PHP_EOL;

    // メッセージ本体
    $messageBody = json_encode([
        'message' => 'Hello from PHP!',
        'timestamp' => time(),
    ]);

    echo 'Message Body: ' . $messageBody . PHP_EOL;

    // 送信するメッセージ属性を `_datadog` に JSON 形式で含める
    $messageAttributes = [
        '_datadog' => [
            'DataType' => 'String',
            'StringValue' => json_encode([
                'traceparent' => $traceparent
            ]),
        ],
    ];

    // メッセージ属性を出力
    echo 'Message Attributes: ' . json_encode($messageAttributes, JSON_PRETTY_PRINT) . PHP_EOL;

    // メッセージを送信
    $result = $sqsClient->sendMessage([
        'QueueUrl'    => $queueUrl,
        'MessageBody' => $messageBody,
        'MessageAttributes' => $messageAttributes,
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
