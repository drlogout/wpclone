FROM wordpress:latest

RUN apt-get update && apt-get install -y \
    default-mysql-client \
    less \
    openssh-client \
    rsync \
    vim

WORKDIR /var/www/html

# install wp cli to /usr/local/bin
RUN curl -o /usr/local/bin/wp.phar https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar && \
    chmod +x /usr/local/bin/wp.phar

COPY wp.sh /usr/local/bin/wp
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
