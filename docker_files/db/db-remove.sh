#!/bin/bash

set -eo pipefail

# Variables
DB_NAME=$1
DB_USER=$1
REMOTE_HOST=%

export MYSQL_PWD="$MARIADB_ROOT_PASSWORD" 

# Check if the database and user exist
set +e
DB_EXISTS=$(mariadb -uroot -sse "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '$DB_NAME';")
USER_EXISTS=$(mariadb -uroot -sse "SELECT EXISTS(SELECT 1 FROM mysql.user WHERE user = '$DB_USER');")
set -e

if [ "$USER_EXISTS" -ne 0 ]; then
    mariadb -uroot -e "DROP USER '$DB_USER'@'$REMOTE_HOST';"
    printf "User %s removed successfully.\n" "$DB_USER"
fi

if [ "$DB_EXISTS" ]; then
    mariadb -uroot -e "DROP DATABASE $DB_NAME;"
    printf "Database %s removed successfully.\n" "$DB_NAME"
fi
