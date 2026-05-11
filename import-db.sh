#!/bin/bash

# ================================================
#   Import SQL ke MySQL Docker
# ================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()   { echo -e "${GREEN}[OK]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# --- Cek file .env ---
if [ ! -f ".env" ]; then
  error "File .env tidak ditemukan!"
fi

# --- Load env ---
set -o allexport
source .env
set +o allexport

# --- Validasi parameter ---
if [ -z "$1" ]; then
  error "Gunakan: ./import-db.sh namafile.sql"
fi

SQL_FILE=$1

if [ ! -f "$SQL_FILE" ]; then
  error "File $SQL_FILE tidak ditemukan!"
fi

# --- Cek container MySQL ---
if ! docker ps --format '{{.Names}}' | grep -q "db_chicken_shop"; then
  error "Container db_chicken_shop tidak berjalan!"
fi

# --- Import ---
echo ""
warn "Mengimport $SQL_FILE ke database $DB_NAME ..."
echo ""

docker exec -i db_chicken_shop \
  mysql -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < "$SQL_FILE"

if [ $? -eq 0 ]; then
  log "Import database berhasil!"
else
  error "Import gagal!"
fi
