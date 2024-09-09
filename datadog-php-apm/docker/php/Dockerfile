FROM php:8.1.18-fpm
COPY php.ini /usr/local/etc/php/

RUN apt-get update && apt-get install -y \
	git \
	curl \
	zip \
	unzip \
    && docker-php-ext-install pdo_mysql \
    && docker-php-ext-install bcmath

# Composerのインストール
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

WORKDIR /var/www

COPY composer.json composer.lock /var/www/
RUN /usr/bin/composer install --no-interaction --no-plugins --no-scripts
RUN /usr/bin/composer require monolog/monolog
RUN /usr/bin/composer require aws/aws-sdk-php


# Datadog Tracerのインストール
COPY datadog-setup.php .
RUN php datadog-setup.php --php-bin=all