#!/bin/bash

# Load .env file
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".env file not found!"
    exit 1
fi

# Check if MONGO_DB_ATLAS_CREDENTIALS is set
if [ -z "$MONGO_DB_ATLAS_CREDENTIALS" ]; then
    echo "MONGO_DB_ATLAS_CREDENTIALS is not set in the .env file!"
    exit 1
fi

# Import each JSON file in the catalogs directory
for file in data/catalogs/*.json; do
  echo "Importing $file..."
  mongoimport --uri "$MONGO_DB_ATLAS_CREDENTIALS" \
              --db adviseu_db \
              --collection catalogs \
              --file "$file" \
              --jsonArray \
              --upsertFields "course_number"
  if [ $? -ne 0 ]; then
    echo "Failed to import $file"
    exit 1
  fi
done

echo "All files imported successfully."
