#!/usr/bin/env bash
set -Eeuo pipefail

SSH_PUBLIC_KEY=${SSH_PUBLIC_KEY:-}

screen -dmS sshd /usr/sbin/sshd -D

# Allow www-data to login
usermod --shell /bin/bash www-data
chown -R www-data:www-data /var/www

su www-data -c /usr/local/bin/install-wordpress.sh

# if SSH_PUBLIC_KEY not empty
if [ -n "$SSH_PUBLIC_KEY" ]; then
    echo "Adding SSH public key"
    mkdir -p /var/www/.ssh
    echo "$SSH_PUBLIC_KEY" > /var/www/.ssh/authorized_keys
    chown -R www-data:www-data /var/www/.ssh
    chmod 700 /var/www/.ssh
    chmod 600 /var/www/.ssh/authorized_keys
fi

exec "$@"
