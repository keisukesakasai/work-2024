<?php

require '../vendor/autoload.php';

use Monolog\Logger;
use Monolog\Handler\StreamHandler;
use Monolog\Formatter\JsonFormatter;
use DDTrace\GlobalTracer;

// AWS 署名バージョン 4 に必要な関数
function sign($key, $msg) {
    return hash_hmac('sha256', $msg, $key, true);
}

function getSignatureKey($key, $dateStamp, $regionName, $serviceName) {
    $kDate = sign('AWS4' . $key, $dateStamp);
    $kRegion = sign($kDate, $regionName);
    $kService = sign($kRegion, $serviceName);
    $kSigning = sign($kService, 'aws4_request');
    return $kSigning;
}

$awsKey = $_ENV['AWS_KEY'];
$awsSecret = $_ENV['AWS_SECRET'];
$queueUrl = 'https://sqs.ap-northeast-1.amazonaws.com/601427279990/sakasai-work-sqs'; // SQS キュー URL
$region = 'ap-northeast-1';
$service = 'sqs';

// メッセージ本体
$messageBody = 'Hello from PHP with HTTP request!';

// 日時の設定
$now = new DateTime('UTC');
$amzDate = $now->format('Ymd\THis\Z');
$dateStamp = $now->format('Ymd');

// Datadog トレーサーの取得
$tracer = GlobalTracer::get();
$scope = $tracer->startActiveSpan('curl_request_span');
$span = $scope->getSpan();

// Datadog のトレース ID を 16進数から X-Amzn-Trace-Id の形式に変換
$traceIdHex = $span->getContext()->getTraceId();  // 64桁の16進数形式
$traceIdXray = sprintf('1-%s-%s', substr($traceIdHex, 0, 8), substr($traceIdHex, 8, 24));  // X-Ray 用に変換
$spanId = $span->getContext()->getSpanId();  // Self に相当

// AWS X-Ray のヘッダー形式に変換 (X-Amzn-Trace-Id)
$amznTraceHeader = sprintf('Root=%s;Parent=%s;Sampled=1', $traceIdXray, $spanId);

// リクエスト情報
$method = 'POST';
$canonicalUri = '/'; // SQS のエンドポイント
$canonicalQuerystring = ''; // クエリストリングなし

$canonicalHeaders = "content-type:application/x-www-form-urlencoded\nhost:sqs.{$region}.amazonaws.com\nx-amz-date:{$amzDate}\n";
$signedHeaders = 'content-type;host;x-amz-date';

$requestParameters = http_build_query([
    'Action' => 'SendMessage',
    'MessageBody' => $messageBody,
    'QueueUrl' => $queueUrl,
    'Version' => '2012-11-05',
]);

// 署名の準備
$payloadHash = hash('sha256', $requestParameters);
$canonicalRequest = "{$method}\n{$canonicalUri}\n{$canonicalQuerystring}\n{$canonicalHeaders}\n{$signedHeaders}\n{$payloadHash}";
$algorithm = 'AWS4-HMAC-SHA256';
$credentialScope = "{$dateStamp}/{$region}/{$service}/aws4_request";
$stringToSign = "{$algorithm}\n{$amzDate}\n{$credentialScope}\n" . hash('sha256', $canonicalRequest);
$signingKey = getSignatureKey($awsSecret, $dateStamp, $region, $service);
$signature = hash_hmac('sha256', $stringToSign, $signingKey);

// Authorization ヘッダーの作成
$authorizationHeader = "{$algorithm} Credential={$awsKey}/{$credentialScope}, SignedHeaders={$signedHeaders}, Signature={$signature}";

// cURL リクエストの実行
$ch = curl_init("https://sqs.{$region}.amazonaws.com/");
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_HTTPHEADER, [
    "Content-Type: application/x-www-form-urlencoded",
    "X-Amz-Date: {$amzDate}",
    "Authorization: {$authorizationHeader}",
    "X-Amzn-Trace-Id: {$amznTraceHeader}",  // X-Amzn-Trace-Id ヘッダーにトレース情報を追加
]);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, $requestParameters);

$response = curl_exec($ch);
if ($response === false) {
    echo 'Curl error: ' . curl_error($ch);
} else {
    echo 'Response: ' . $response;
}
curl_close($ch);

// スパンを終了
$span->finish();
