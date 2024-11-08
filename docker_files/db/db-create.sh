#!/bin/bash

set -eo pipefail

# Variables
DB_NAME=$1
DB_USER=$1
DB_PASS=$1
REMOTE_HOST=%
CHARACTER_SET='utf8mb4'
COLLATION='utf8mb4_unicode_ci'

export MYSQL_PWD="$MARIADB_ROOT_PASSWORD" 

# Check if the database and user exist
set +e
DB_EXISTS=$(mariadb -uroot -sse "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '$DB_NAME';")
USER_EXISTS=$(mariadb -uroot -sse "SELECT EXISTS(SELECT 1 FROM mysql.user WHERE user = '$DB_USER');")
set -e

if [ ! "$DB_EXISTS" ]; then
    mariadb -uroot -e "CREATE DATABASE $DB_NAME CHARACTER SET $CHARACTER_SET COLLATE $COLLATION;"
    printf "Database %s created successfully.\n" "$DB_NAME"
fi

if [ "$USER_EXISTS" -eq 0 ]; then
    mariadb -uroot -e "CREATE USER '$DB_USER'@'$REMOTE_HOST' IDENTIFIED BY '$DB_PASS';"

    mariadb -uroot -e "GRANT ALL PRIVILEGES ON $DB_NAME.* TO '$DB_USER'@'$REMOTE_HOST';"
    mariadb -uroot -e "FLUSH PRIVILEGES;"

    printf "User %s created successfully.\n" "$DB_USER"
fi
