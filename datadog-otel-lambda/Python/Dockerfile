# ベースイメージとして Python を使用
FROM python:3.11-slim

# 作業ディレクトリを設定
WORKDIR /app

# 必要な依存関係をインストール
COPY requirements.txt .
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/* \
    && pip install --no-cache-dir -r requirements.txt

# アプリケーションのソースコードをコピー
COPY . .

# 必要な環境変数を設定
ENV SERVER_PORT=80

# コンテナのポートを公開
EXPOSE 80

# アプリケーションを ddtrace-run を使って実行
CMD ["ddtrace-run", "python", "main.py"]
