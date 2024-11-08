#!/usr/bin/env bash
set -Eeuo pipefail

# Allow www-data to login, set home directory and change ownership of /var/www
if [ ! -d /home/www-data ]; then
    mkdir -p /home/www-data
fi
usermod --shell /bin/bash --home /home/www-data www-data

set +e
chown -R www-data:www-data /home/www-data
chown -R www-data:www-data /var/www
set -e

exec "$@"
