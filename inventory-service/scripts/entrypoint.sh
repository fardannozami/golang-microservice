#!/bin/sh
set -e

# Jalankan seeder jika environment variable RUN_SEEDER=true
if [ "$RUN_SEEDER" = "true" ]; then
    echo "Running database seeder..."
    /app/seed
    echo "Seeding completed!"
fi

# Jalankan aplikasi utama
exec /app/inventory-service