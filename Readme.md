# File Sharing & Management System (Backend)

## Project Description

This project is a backend service for a file-sharing platform that allows users to upload, manage, and share files. The system handles multiple users, stores metadata in PostgreSQL, manages file uploads to S3, and implements caching for file metadata. The project is built in Go, demonstrating proficiency in handling concurrency and performance optimizations.

## Features

### 1. User Authentication & Authorization

- **Endpoints**:
  - **POST /register**: Register a new user with email and password.
  - **POST /login**: Log in with email and password to receive a JWT access token.

- **Access Token Usage**:
  - Upon successful login, an `access_token` is returned.
  - This token must be included in the `Authorization` header as `Bearer <token>` for accessing protected routes such as `/files`.

### 2. File Upload & Management

- **Endpoints**:
  - **POST /upload**: Upload files (documents, images) to S3 or local storage.
  - **GET /files**: Retrieve metadata for all uploaded files for the authenticated user.
  - **GET /share/:file_id**: Get a presigned URL for sharing a specific file.
  - **PATCH /update/:file_id**: Update file metadata.
  - **DELETE /delete/:file_id**: Delete a specific file.

- **File Upload**:
  - Files are chunked into 5MB parts and uploaded concurrently to S3 using multipart upload.

- **File Sharing**:
  - Generate a presigned URL for the file using S3, allowing secure access to the file via a public link.

### 3. File Retrieval & Sharing

- **Endpoints**:
  - **GET /files**: Retrieve all files' metadata for the authenticated user.
  - **GET /share/:file_id**: Get a presigned URL for sharing the file.

### 4. File Search

- **Search Parameters**:
  - `file_name`
  - `start_date`
  - `end_date`
  - `file_type`

- **Indexing**:
  - Indexes are created for efficient searching by file name, upload date, and file type.

### 5. Caching Layer for File Metadata

- **Caching**:
  - Metadata is cached on retrieval using Redis.
  - Cache is reset when files are uploaded, updated, or deleted.

### 6. Database Schema

- **Tables**:
  - **users**: Stores user information.
  - **files**: Stores file metadata.
  - **shared_files**: Tracks files shared with users.

- **Indexes**:
  - Indexes on user_id, file_name, upload_date, and file_type for efficient searching.

- **Migrations**:
  - Managed using Atlas/Go for schema migrations and versioning.

### 7. Background Job for File Deletion

- **Scheduled Tasks**:
  - **Shared File Deletion**: Runs every 10 minutes to remove expired shared files.
  - **Expired File Deletion**: Runs twice a day to remove expired files from S3.

- **Tool**:
  - Implemented using the `gocron` library for scheduling tasks.

### 8. Testing

- **Tests**:
  - Comprehensive tests are located in the `tests` directory to validate API functionality.


### 9. Deployment

- **Deployment**:
  - **CI/CD Pipeline**: The project is automatically built and deployed using GitHub Actions. The pipeline is configured to trigger on pushes to the `main` branch. It uses SSH to connect to the server, pulls the latest code from the repository, and manages Docker containers to ensure the application is up-to-date.
  - **Deployment Steps**:
    1. **SSH Connection**: The GitHub Actions workflow connects to the server via SSH.
    2. **Pull Latest Code**: It pulls the latest changes from the `main` branch of the repository.
    3. **Docker Management**: The workflow stops any running Docker containers, rebuilds the Docker images, and restarts the containers with the new build.
  - **Server Details**:
    - The application is hosted on a Microsoft Azure VM.
    - **Live API URL**: [http://20.40.48.251](http://20.40.48.251) provides access to the deployed application.

### 10. Dockerfile

The project uses a two-stage Dockerfile to minimize the final image size:

- **Build Stage**:
  - Uses `golang:1.21.0-alpine3.18` as the base image.
  - Installs dependencies, compiles the Go application, and creates a binary executable.

- **Run Stage**:
  - Uses `alpine:3.18` as the base image.
  - Copies the compiled binary and environment files into the image.
  - Exposes port 8080 and sets the command to run the application.

## Environment Variables

Create a `.env` file in the root directory with the following environment variables:

```env
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=admin
POSTGRES_PASSWORD=password123
POSTGRES_DB=file_management_system

CLIENT_ORIGIN=*
PORT=8080

ACCESS_SECRET_KEY=<your_access_secret_key>
REFRESH_SECRET_KEY=<your_refresh_secret_key>

REDIS_HOST=redis:6379
REDIS_PASSWORD=redis123
REDIS_DB=0

AWS_ACCESS_KEY_ID=<your_aws_access_key_id>
AWS_SECRET_ACCESS_KEY=<your_aws_secret_access_key>
AWS_REGION=ap-south-1
AWS_BUCKET=file-upload-trademarkia-bucket
```

## Running the Project

### Local Development

1. **Clone the Repository**:
   ```sh
   git clone <repository_url>
   cd <repository_directory>
   ```

2. **Environment Variables**:
   Create a `.env` file in the root directory and configure the environment variables as specified above.

3. **Run with Docker Compose**:
   ```sh
   docker compose up -d --build
   ```

## Postman Documentation

Detailed API documentation is available via Postman:
- [API Documentation](https://documenter.getpostman.com/view/21877920/2sAXqp83yu#7c9d7cc7-31c5-44eb-bc54-ed14d2e47887)

---
## Architecture

### Layered Architecture

The project follows a layered architecture to organize code and separate concerns:

1. **Handlers**:
   - Responsible for processing incoming HTTP requests, validating input, and returning responses.

2. **Services**:
   - Contains business logic and interacts with the database and other components.

3. **DB**:
   - Manages database operations and interactions using SQLC for query generation and database schema management.

4. **Middleware**:
   - Provides functionalities such as authentication, logging, and error handling.

5. **Routes**:
   - Defines the API endpoints and maps them to the appropriate handlers.

## Features Checklist

- [x] User Authentication & Authorization
- [x] File Upload & Management
- [x] File Retrieval & Sharing
- [x] File Search
- [x] Caching Layer for File Metadata
- [x] Database Schema
- [x] Background Job for File Deletion
- [x] Testing

### Bonus Tasks (Optional)

- [ ] WebSocket for Real-Time File Upload Notifications
- [ ] File Encryption
- [x] Hosting
- [x] Implement Rate Limiting (100 requests per user per minute)

