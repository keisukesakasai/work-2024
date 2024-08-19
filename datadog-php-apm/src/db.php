<?php

  require '../vendor/autoload.php';

  // load Monolog library
  use Monolog\Logger;
  use Monolog\Handler\StreamHandler;
  use Monolog\Formatter\JsonFormatter;

  $host = 'mysql'; // MySQLコンテナのサービス名
  $dbname = $_ENV['MYSQL_DATABASE'];
  $username = 'root';
  $password = $_ENV['MYSQL_ROOT_PASSWORD'];

  // ログチャネルを作成
  $log = new Logger('php-fpm');

  // Json フォーマッタを作成
  $formatter = new JsonFormatter();

  // ハンドラーを作成
  $stream = new StreamHandler('/log/application-json.log', Logger::DEBUG);
  $stream->setFormatter($formatter);

  // バインド
  $log->pushHandler($stream);

try {
    // データベース接続のログ
    $log->info("Connecting to the database...");

    // 新しいPDOオブジェクトを作成し、MySQLデータベースに接続
    $db = new PDO("mysql:host={$host};dbname={$dbname};charset=utf8", $username, $password);

    // PDOのエラーモードを例外に設定
    $db->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    $log->info("Connection successful!");

    // SQL文を実行
    $log->info("Executing SQL query...");
    $stmt = $db->prepare('SELECT * FROM mytable');
    $stmt->execute();

    $results = $stmt->fetchAll(PDO::FETCH_ASSOC);
    $log->info("Query executed successfully. Retrieved " . count($results) . " rows.");

    foreach ($results as $result) {
        echo $result['id'] . '. ' . $result['name'] . PHP_EOL;
    }

} catch (PDOException $e) {
    // PDOのエラーをログに記録し、標準出力にも表示
    $errorMessage = 'Database Error: ' . $e->getMessage();
    $log->error($errorMessage);
    echo $errorMessage . PHP_EOL;
    echo 'An error occurred while accessing the database. Please check the logs for details.' . PHP_EOL;
} catch (Exception $e) {
    // その他のエラーをログに記録し、標準出力にも表示
    $errorMessage = 'General Error: ' . $e->getMessage();
    $log->error($errorMessage);
    echo $errorMessage . PHP_EOL;
    echo 'An unexpected error occurred. Please check the logs for details.' . PHP_EOL;
}
