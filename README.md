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
    Use Docker Compose to build and run the application and the database. This setup includes `air` for automatic hot-reloading when you change a file.

    ```bash
    docker-compose up --build
    ```
    The first time you run this, it will install `air`. Subsequent runs will be faster.

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

    c. **Install `air`:** Install the `air` tool on your local machine.
    ```bash
    go install github.com/cosmtrek/air@latest
    ```

    d. **Run Migrations & Start the App:** Run the migrations and then start the application using `air`.
    ```bash
    # First, run migrations (only needed once or when schema changes)
    migrate -path migrations -database "postgres://user:password@localhost:5432/blog?sslmode=disable" up

    # Then, run the app with hot-reloading
    air
    ```

The application will be available at `http://localhost:8080`.