# AdviseU backend

Notes on the backend.

Dependencies:
```bash
go install github.com/air-verse/air@latest
go get go.mongodb.org/mongo-driver/mongo
go get github.com/joho/godotenv
```

Run the server with automatic restarts whenever Go files are changed:
```bash
air
```

Database:
MongoDB Atlas is being used to host the cloud databse. To access this database, create a `.env` file containing your MongoDB credentials which start with `mongodb+srv://`: `MONGO_DB_ATLAS_CREDENTIALS`.
