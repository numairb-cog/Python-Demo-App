# Go Gin Demo App

A demo web application built with Go and the Gin framework that demonstrates various features for performance monitoring and testing. This app has routes that:

1. Sleep in a way that draws a pretty response time graph (sine wave pattern)
2. Raise uncaught exceptions in application code
3. Simulate slow database calls in PostgreSQL
4. Make external HTTP calls
5. Return HTTP errors (5xx and 4xx)

## Requirements

- Go 1.18 or higher
- PostgreSQL database (optional - app will run without it but /query routes won't work)

## Installation

Clone the repository:

```bash
git clone https://github.com/numairb-cog/Python-Demo-App.git
cd Python-Demo-App
```

Install dependencies:

```bash
go mod download
```

## Configuration

The application uses environment variables for configuration:

- `DEMO_PGSQL_USER` - PostgreSQL username (default: "test")
- `DEMO_PGSQL_PASSWORD` - PostgreSQL password (default: "test")
- `DEMO_PGSQL_HOST` - PostgreSQL host (default: "127.0.0.1")
- `DEMO_PGSQL_DB` - PostgreSQL database name (default: "test")
- `PORT` - Server port (default: "9000")

## Database Setup (Optional)

If you want to use the database query routes, install and configure PostgreSQL:

```bash
# The database configuration defaults to:
# - Host: 127.0.0.1
# - Port: 5432 (default PostgreSQL port)
# - User: test
# - Password: test
# - Database: test

# Create a test database and user (example for PostgreSQL):
sudo -u postgres psql -c "CREATE USER test WITH PASSWORD 'test';"
sudo -u postgres psql -c "CREATE DATABASE test OWNER test;"
```

You don't need any tables in the database - the app just executes simple queries for testing.

## Running the Application

### Development Mode

```bash
go run main.go
```

### Production Build

```bash
go build -o demo-app
./demo-app
```

The web server runs on port 9000 by default. Access it at http://localhost:9000

## Available Routes

- `/` - Index page with links to all demo routes
- `/wave/<id>` - Variable response time with sine wave delay (e.g., `/wave/123`)
- `/error/always` - Always raises an exception
- `/error/sometimes` - Raises an exception 10% of the time
- `/query/pgsql` - PostgreSQL query simulator (slow/error/normal queries)
- `/http?url=<url>` - HTTP exit call (e.g., `/http?url=http://example.com`)

## Testing

Build and test the application:

```bash
# Format code
go fmt ./...

# Check for common errors
go vet ./...

# Build the application
go build

# Run the application
go run main.go
```

Then visit http://localhost:9000 in your browser to test all routes.

## Generating Load

You can use tools like `siege`, `ab` (Apache Bench), or `wrk` to generate load:

```bash
# Using siege
siege -d 1 -c 10 http://localhost:9000/

# Using Apache Bench
ab -n 1000 -c 10 http://localhost:9000/wave/1

# Using wrk
wrk -t10 -c100 -d30s http://localhost:9000/
```

## Dependencies

- [Gin Web Framework](https://github.com/gin-gonic/gin) - HTTP web framework
- [pq](https://github.com/lib/pq) - PostgreSQL driver

## License

See LICENSE file for details.
