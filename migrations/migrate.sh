#!/bin/bash

DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-"victor12"}
DB_NAME=${DB_NAME:-"postgres"}
MIGRATIONS_DIR=${MIGRATIONS_DIR:-"./migrations"}

DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
echo "Using db url: $DB_URL"

if ! command -v migrate &> /dev/null; then
    echo "Error: migrate tool is not installed"
    echo "Please install it with: brew install golang-migrate"
    exit 1
fi

echo "Waiting for database to be ready..."
while ! pg_isready -h "$DB_HOST" -p "$DB_PORT"; do
    sleep 1
done
echo "Database is ready."

run_migration() {
    case $1 in
        "up")
            echo "Running migrations up..."
            migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" -verbose up || { echo "Migration failed"; exit 1; }
            ;;
        "down")
            migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" -verbose down || { echo "Migration failed"; exit 1; }
            ;;
        "steps")
            if [ -z "$2" ]; then
                echo "Please provide the number of steps."
                exit 1
            fi
            migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" "$2" || { echo "Migration failed"; exit 1; }
            ;;
        *)
            echo "Invalid option. Use 'up', 'down', or 'steps <n>'."
            exit 1
            ;;
    esac
}

run_migration "$@"