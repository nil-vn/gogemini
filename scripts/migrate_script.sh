#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

USERNAME=""
PASSWORD=""
DB_URL="${DB_URL:-sqlite3://app.db}"
DB_DRIVER="${DB_DRIVER:-sqlite}"
DB_DSN="${DB_DSN:-file:app.db?cache=shared}"
ROLE="${ROLE:-admin}"
STATUS="${STATUS:-active}"
EMAIL="${EMAIL:-}"

usage() {
  cat <<EOF
Usage: scripts/migrate_script.sh --username=<admin> --password=<password> [options]
Options:
  --db-url=<url>        Migration DB URL (default: ${DB_URL})
  --db-driver=<driver>  SQL driver for bootstrap insert (default: ${DB_DRIVER})
  --db-dsn=<dsn>        SQL DSN for bootstrap insert (default: ${DB_DSN})
  --role=<role>         Default user role (default: ${ROLE})
  --status=<status>     Default user status (default: ${STATUS})
  --email=<email>       Optional email (default: <username>@local)
EOF
}

for arg in "$@"; do
  case "$arg" in
    --username=*) USERNAME="${arg#*=}" ;;
    --password=*) PASSWORD="${arg#*=}" ;;
    --db-url=*) DB_URL="${arg#*=}" ;;
    --db-driver=*) DB_DRIVER="${arg#*=}" ;;
    --db-dsn=*) DB_DSN="${arg#*=}" ;;
    --role=*) ROLE="${arg#*=}" ;;
    --status=*) STATUS="${arg#*=}" ;;
    --email=*) EMAIL="${arg#*=}" ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown argument: $arg"; usage; exit 1 ;;
  esac
done

if [[ -z "$USERNAME" || -z "$PASSWORD" ]]; then
  echo "--username and --password are required"
  usage
  exit 1
fi

if ! command -v migrate >/dev/null 2>&1; then
  echo "golang-migrate CLI is required: https://github.com/golang-migrate/migrate"
  exit 1
fi

if [[ -z "$EMAIL" ]]; then
  EMAIL="${USERNAME}@local"
fi

echo "[1/2] Applying migrations to ${DB_URL}"
migrate -path migrations -database "$DB_URL" up

echo "[2/2] Ensuring default admin user '${USERNAME}' exists"
go run ./cmd/bootstrap-admin \
  --db-driver="$DB_DRIVER" \
  --db-dsn="$DB_DSN" \
  --username="$USERNAME" \
  --password="$PASSWORD" \
  --role="$ROLE" \
  --status="$STATUS" \
  --email="$EMAIL"

echo "Done"
