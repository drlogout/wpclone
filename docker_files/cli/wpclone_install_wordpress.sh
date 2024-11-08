#!/usr/bin/env bash
set -Eeuo pipefail

DEFAULT_TITLE="My wpclone Site"
DEFAULT_LOCALE="en_US"
DEFAULT_ADMIN_USER="admin"
DEFAULT_ADMIN_EMAIL="admin@example.com"

while [[ "$#" -gt 0 ]]; do
    case $1 in
        --db-name) DB_NAME="$2"; shift ;;
        --db-host) DB_HOST="$2"; shift ;;
        --db-user) DB_USER="$2"; shift ;;
        --db-password) DB_PASSWORD="$2"; shift ;;
        --url) URL="$2"; shift ;;
        --title) TITLE="$2"; shift ;;
        --admin-user) ADMIN_USER="$2"; shift ;;
        --admin-password) ADMIN_PASSWORD="$2"; shift ;;
        --admin-email) ADMIN_EMAIL="$2"; shift ;;
        --locale) LOCALE="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

function is_installed() {
    if [ ! -f "./wp-config.php" ]; then
        return 1
    fi

    return 0
}

# sleep until mysql is ready 
timeout=60
until mysqladmin ping -h "${DB_HOST}" --silent; do
    sleep 1
    ((timeout--))
    if [ $timeout -le 0 ]; then
        echo "MySQL is not ready after waiting for 60 seconds."
        exit 1
    fi
done

if ! is_installed; then
    echo "Installing WordPress..."
    wp core download \
        --locale="${LOCALE:-$DEFAULT_LOCALE}"

    wp config create \
        --dbhost="${DB_HOST}" \
        --dbname="${DB_NAME}" \
        --dbuser="${DB_USER}" \
        --dbpass="${DB_PASSWORD}"

    wp core install \
        --url="${URL}" \
        --title="${TITLE:-$DEFAULT_TITLE}" \
        --admin_user="${ADMIN_USER:-$DEFAULT_ADMIN_USER}" \
        --admin_password="${ADMIN_PASSWORD:-$DEFAULT_ADMIN_USER}" \
        --admin_email="${ADMIN_EMAIL:-$DEFAULT_ADMIN_EMAIL}"
fi
