#!/usr/bin/env bash

export MYSQL_PWD="$MARIADB_ROOT_PASSWORD"

# Database connection parameters
DB_HOST=$1
DB_PORT=$1
DB_USER=$1
DB_NAME=$1

# Timeout settings
TIMEOUT=60       # Timeout in seconds
SLEEP_INTERVAL=5 # Interval between connection attempts

# Function to check database connection
check_database_exists() {
    # Query to check if the database exists
    mariadb -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME;" >/dev/null 2>&1
    return $?
}

# Start time
START_TIME=$(date +%s)

# Loop until connection is successful or timeout is reached
while true; do
    if check_connection; then
        echo "Successfully connected to the database."
        exit 0
    fi

    CURRENT_TIME=$(date +%s)
    ELAPSED_TIME=$((CURRENT_TIME - START_TIME))

    if [ "$ELAPSED_TIME" -ge "$TIMEOUT" ]; then
        echo "Failed to connect to the database within $TIMEOUT seconds."
        exit 1
    fi

    echo "Waiting for database connection..."
    sleep "$SLEEP_INTERVAL"
done
