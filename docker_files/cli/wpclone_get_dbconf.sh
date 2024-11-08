#!/usr/bin/env bash
CONFIG_FILE="wp-config.php"
if [ -n "$1" ]; then
    CONFIG_FILE=$1
fi

DB_HOST=$(wp config get DB_HOST --config-file="$CONFIG_FILE")
DB_NAME=$(wp config get DB_NAME --config-file="$CONFIG_FILE")
DB_USER=$(wp config get DB_USER --config-file="$CONFIG_FILE")
DB_PASSWORD=$(wp config get DB_PASSWORD --config-file="$CONFIG_FILE")

echo "{\"host\": \"$DB_HOST\", \"name\": \"$DB_NAME\", \"user\": \"$DB_USER\", \"password\": \"$DB_PASSWORD\"}"
