# Go Gin Benchmarks Example

[![Go Report Card](https://goreportcard.com/badge/github.com/CodSpeedHQ/go-gin-benchmarks-example)](https://goreportcard.com/report/github.com/CodSpeedHQ/go-gin-benchmarks-example)
[![CodSpeed Badge](https://img.shields.io/endpoint?url=https://codspeed.io/badge.json)](https://codspeed.io/CodSpeedHQ/go-gin-benchmarks-example)

Example repository demonstrating how to benchmark a Gin HTTP API with CodSpeed.

## Running the API

```bash
# Install dependencies
go mod tidy

# Run the server
go run api.go
```

The API will be available at `http://localhost:8080` with these endpoints:
- `GET /albums` - List all albums
- `GET /albums/:id` - Get album by ID
- `POST /albums` - Create a new album

## Running Benchmarks

```bash
# Run all benchmarks
go test -bench=.

# Run with longer benchmark time for more accurate results
go test -bench=. -benchtime=5s
```

## CI Workflow

The repository uses GitHub Actions with CodSpeed to:
- Run benchmarks on every pull request using `codspeed-macro` runners
- Track performance changes over time
- Automatically comment on PRs with performance impact

See the [full guide](https://docs.codspeed.io/guides/benchmarking-a-go-gin-api) for implementation details.
