# Simple HTTP session checker

## About <a name = "about"></a>
Simple HTTP session (cookie, redis) checker. By GET /incr, you will get the following response.

```
{"count":3,"hostname":"cf238ff21641"}
```

Count will increase each time you GET /incr depending on the session.

## Installation

1. Clone the repository:

```sh
$ git clone https://github.com/torumakabe/session-checker.git
$ cd session-checker
```

2. Install dependencies:

```sh
$ go mod tidy
```

## Usage

1. Build the application:

```sh
$ make build
```

2. Run the application:

```sh
$ ./bin/session-checker
```

## Testing

Run the tests:

```sh
$ make test
```

## Configuration

The application can be configured using environment variables or command-line flags:

- `--redis-server` or `-r`: Set Redis server hostname:port.
- `--redis-password` or `-p`: Set Redis password.

Environment variables:

- `SESSION_CHECKER_REDIS_SERVER`: Set Redis server hostname:port.
- `SESSION_CHECKER_REDIS_PASSWORD`: Set Redis password.
