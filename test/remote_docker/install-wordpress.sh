#!/usr/bin/env bash
set -Eeuo pipefail

DEFAULT_WORDPRESS_TITLE="My WordPress Site"
DEFAULT_WORDPRESS_LOCALE="en_US"

function is_wordpress_installed() {
    if [ ! -f "./wp-config.php" ]; then
        return 1
    fi

    return 0
}

# sleep until mysql is ready
until mysqladmin ping -h "${WORDPRESS_DB_HOST}" --silent; do
    sleep 1
done

if ! is_wordpress_installed; then
    echo "Installing WordPress..."

    wp core download \
        --locale="${WORDPRESS_LOCALE:-$DEFAULT_WORDPRESS_LOCALE}"

    wp config create \
        --dbhost="${WORDPRESS_DB_HOST}" \
        --dbname="${WORDPRESS_DB_NAME}" \
        --dbuser="${WORDPRESS_DB_USER}" \
        --dbpass="${WORDPRESS_DB_PASSWORD}"

    wp core install \
        --url="${WORDPRESS_URL}" \
        --title="${WORDPRESS_TITLE:-$DEFAULT_WORDPRESS_TITLE}" \
        --admin_user="${WORDPRESS_ADMIN_USER}" \
        --admin_password="${WORDPRESS_ADMIN_PASSWORD}" \
        --admin_email="${WORDPRESS_ADMIN_EMAIL}"
        
    rm -rf wp-content/themes/twenty*

    mv /tmp/blankslate wp-content/themes/blankslate

    wp theme activate blankslate
fi
