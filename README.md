# Go Blog Application

This is a blog application built with Go, Echo, PostgreSQL, and deployed on AWS.

## Requirements

- Go 1.21+
- Docker & Docker Compose
- An AWS account (for S3 storage in production)

## Local Development

1.  **Choose Your Method:**

    -   **With Docker (Recommended):** The easiest way to get started. It runs the application and a PostgreSQL database.
    -   **Without Docker:** Useful for faster iteration during development. You will need to run your own PostgreSQL instance.

2.  **Run with Docker:**
    Use Docker Compose to build and run the application and the database.

    ```bash
    docker-compose up --build
    ```

3.  **Run Migrations:**
    You will need a migration tool like `golang-migrate/migrate` to run the SQL migrations against the database.

    ```bash
    # For the Docker container DB
    migrate -path migrations -database "postgres://user:password@localhost:5432/blog?sslmode=disable" up
    ```

4.  **Run without Docker:**

    a. **Install PostgreSQL:** Make sure you have a PostgreSQL server running locally.

    b. **Set Environment Variables:** Copy the `.env.example` file to a new file named `.env`.
    ```bash
    cp .env.example .env
    ```
    Update the `DATABASE_URL` in `.env` if your local PostgreSQL setup is different.

    c. **Run Migrations:** Run the migrations against your local database.

    d. **Run the Application:**
    ```bash
    go run ./cmd/server/main.go
    ```

The application will be available at `http://localhost:8080`.