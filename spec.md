
Requirements
as engineering manager, manage the team using go and echo framework, postgresSQL and aws. Now we start the project that makes the blog application which allow user to write the post in md format. Everything should be deployed on aws. Let tailor the project code base, which following requirements:

1. go lang and echo framework. The frontend can use bootstrap as well

2. user registestration, create, update their post. There is an pulic page that we can see all the post with pagination. 

3. post in md format, store as file on local test, but we should configure to store the file on s3. Each user has their own s3 bucket

4. We want to track the changes and history on each post

5. The project structure
├── cmd/server/
│   └── main.go                 # Main application entry point
├── internal/
│   ├── api/                    # Echo handlers & routing
│   │   ├── post_handler.go
│   │   ├── user_handler.go
│   │   └── routes.go
│   ├── config/
│   │   └── config.go             # Configuration loading (Viper)
│   ├── middleware/
│   │   └── auth.go               # Auth middleware (e.g., JWT)
│   ├── model/
│   │   ├── post.go
│   │   └── user.go
│   ├── service/                # Business logic
│   │   ├── post_service.go
│   │   └── user_service.go
│   ├── store/                  # Database interfaces
│   │   ├── postgres/           # Postgres implementations
│   │   │   ├── post_store.go
│   │   │   └── user_store.go
│   │   └── store.go
│   └── storage/                # File storage abstraction (Local vs S3)
│       ├── local.go
│       ├── s3.go
│       └── storage.go
├── migrations/                 # SQL database migrations
│   └── 001_init_schema.up.sql
├── .gitignore
├── Dockerfile                  # For our AWS deployment
├── docker-compose.yml          # For local development (app + postgres)
├── go.mod
└── README.md                   # Project overview and setup

6. SSR (Server-Side Rendering): Our main Go (Echo) app fetches the data, renders the Markdown to HTML, and serves a complete web page.