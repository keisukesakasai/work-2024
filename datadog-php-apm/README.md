- Require in `.env` file
  - MYSQL_ROOT_PASSWORD
  - MYSQL_DATABASE
  - MYSQL_USER
  - MYSQL_PASSWORD
  - DD_API_KEY

- Usage
```sh
$ docker compose -f compose.yaml up --build
```

- API
  - `db.php`
  <img src="img/db.png" width=auto height="300">
  
  - `sqs_api_dd.php` ※ Require: SQS / Lambda
  <img src="img/sqs.png" width=auto height="300">