FROM wordpress:latest

RUN apt-get update && apt-get install -y \
    default-mysql-client \
    openssh-server \
    rsync \
    screen \
    vim

# setup sshd
RUN mkdir /var/run/sshd
COPY ssh/sshd_config /etc/ssh/sshd_config
EXPOSE 22

COPY blankslate /tmp/blankslate
RUN chown -R www-data:www-data /tmp/blankslate

# install wp cli to /usr/local/bin
RUN curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar && \
    chmod +x wp-cli.phar && \
    mv wp-cli.phar /usr/local/bin/wp

COPY install-wordpress.sh /usr/local/bin/install-wordpress.sh
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
