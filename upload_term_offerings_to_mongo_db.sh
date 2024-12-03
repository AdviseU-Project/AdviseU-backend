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

# Import each JSON file in the term_offerings directory
for file in data/term_offerings/*.json; do
    echo "Importing $file..."
    # Using 'department' and 'term' as unique identifiers for the term offerings
    mongoimport --uri "$MONGO_DB_ATLAS_CREDENTIALS" \
                --db adviseu_db \
                --collection term_offerings \
                --file "$file" \
                --jsonArray \
                --upsertFields "department,term"
    if [ $? -ne 0 ]; then
        echo "Failed to import $file"
        exit 1
    fi
done

echo "All files imported successfully."
