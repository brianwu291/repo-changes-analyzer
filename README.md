# GitHub Repository Changes Analyzer

A service that analyzes GitHub repository changes, providing insights into code contributions by different users within a specified time range. The service uses goroutines for concurrent processing to enhance performance.

## Project Structure

```
.
├── LICENSE
├── development-container
│   └── Dockerfile
├── config
│   └── config.go
├── go.mod
├── go.sum
├── internal
│   ├── di
│   │   └── di.go
│   ├── handlers
│   │   └── analysishandler
│   │       └── analysishandler.go
│   ├── models
│   │   └── models.go
│   ├── repos
│   │   └── githubrepository
│   │       └── githubrepository.go
│   ├── server
│   │   └── server.go
│   └── services
│       └── analyzerservice
│           └── analyzerservice.go
└── main.go
```

### Component Description

- **config**: Manages application configuration, including GitHub token and server port settings
- **di**: Handles dependency injection, wiring up all components
- **handlers**: Contains HTTP request handlers that process incoming API requests
- **models**: Defines data structures used throughout the application
- **repos**: Implements data access layer, specifically GitHub API interactions
- **server**: Sets up the HTTP server and routing using Gin framework
- **services**: Contains business logic for analyzing repository changes

## Technical Stack

- Golang
- Gin web framework
- Environment variable configuration

## Setup

Can either install directly or through Docker

### Install directly:

1. Clone repository:

```bash
git clone https://github.com/brianwu291/repo-changes-analyzer.git
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
export GITHUB_TOKEN=your_github_token_here
export PORT=8088  # optional, default is 8080
```

or using `.env` file. The dotenv pkg will auto load the env values

4. Run the service:
   `$ go run main.go` or `$ air` for hot reload (need install air first)

### With Docker:

1. Go to `/development-container` folder
2. Run the whole application:
   `docker compose up --build`
   (This include `air` command, so hot-reload supported)

## API Documentation

### Analyze Repository Changes

Analyzes the code changes in a GitHub repository within a specified time range.

**Endpoint**: `POST /api/analyze`

**Request Headers**:

```
Content-Type: application/json
```

**Request Body**:

```json
{
  "owner": "string", // gitHub repo owner
  "repo": "string", // repo name
  "start_date": "string", // start date in YYYY-MM-DD format
  "end_date": "string" // end date in YYYY-MM-DD format
}
```

**Success Response** (200 OK):

```json
{
  "repository": "owner/repo",
  "time_range": "start_date to end_date",
  "user_changes": {
    "username": {
      "additions": 100,
      "deletions": 50,
      "total": 150
    }
    // ... more users
  }
}
```

**Error Response** (400 Bad Request):

```json
{
  "repository": "",
  "time_range": "",
  "user_changes": null,
  "error": "error message describing what went wrong"
}
```

### Example Usage

```bash
curl -X POST http://localhost:8080/api/analyze \
-H "Content-Type: application/json" \
-d '{
    "owner": "brianwu291",
    "repo": "repo-changes-analyzer",
    "start_date": "2025-01-01",
    "end_date": "2025-01-31"
}'
```

<img width="949" alt="截圖 2025-01-03 16 07 13" src="https://github.com/user-attachments/assets/f811c6ed-77e8-4baf-a13e-56bad5ace0f0" />


## Features

- concurrent processing of commits using goroutines
- rate limiting to respect GitHub API limits
- CORS support for web clients

## Contributing

Feel free to open issues or submit pull requests for improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
