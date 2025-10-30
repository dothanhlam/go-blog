# Go Blog Application

This is a blog application built with Go, Echo, PostgreSQL, and deployed on AWS.

## Requirements

- Go 1.21+
- Docker & Docker Compose
- An AWS account (for S3 storage in production)

## Local Development

This project is configured to use `air` for live-reloading, which automatically rebuilds and restarts the application when you save a file. You can run the development environment with or without Docker.

### 1. Run with Docker (Recommended)

The `docker-compose.yml` file is set up to run the Go application and a PostgreSQL database. The app container will automatically install and run `air`.

a. **Start the services:**
```bash
docker-compose up --build
```

b. **Run database migrations:**
In a separate terminal, run the migrations against the Docker database. You only need to do this once or when the database schema changes.
```bash
# For the Docker container DB
migrate -path migrations -database "postgres://user:password@localhost:5432/blog?sslmode=disable" up
```

### 2. Run without Docker

This method gives you faster build times but requires you to manage PostgreSQL and other dependencies on your local machine.

a. **Install Dependencies:**
- Make sure you have a PostgreSQL server running locally.
- Install the `air` live-reloading tool:
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

b. **Set Environment Variables:**
Copy the `.env.example` file to `.env` and update the `DATABASE_URL` if your local PostgreSQL setup is different.
```bash
cp .env.example .env
```

c. **Run Migrations & Start the App:**
First, run the database migrations. Then, start the application using `air`.
```bash
# Run migrations
migrate -path migrations -database "postgres://user:password@localhost:5432/blog?sslmode=disable" up

# Run the app with hot-reloading
air
```

The application will be available at `http://localhost:8080`.