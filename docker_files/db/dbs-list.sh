#!/usr/bin/env bash

export MYSQL_PWD="$MARIADB_ROOT_PASSWORD" 

mariadb -uroot --skip-column-names -s -e "SHOW DATABASES;"
