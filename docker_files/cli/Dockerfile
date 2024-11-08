FROM wordpress:cli-php8.3

COPY wpclone_get_dbconf.sh /usr/local/bin/wpclone_get_dbconf
COPY wpclone_install_wordpress.sh /usr/local/bin/wpclone_install_wordpress
COPY ensure-line.sh /usr/local/bin/ensure-line
COPY remove-line.sh /usr/local/bin/remove-line

WORKDIR /var/www/html

ENTRYPOINT ["docker-entrypoint.sh"]
