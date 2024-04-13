# Poseidon
Poseidon is an open-source project designed to track predefined Uniswap V3 pools in a straightforward, scalable manner, with reorg handled.

## Demo

![Demo GIF]()

## Tech Stack
- [Go](https://go.dev/) - Language
- [Gin](https://gin-gonic.com/) - Web Framework
- [Gorm](https://gorm.io/index.html) - ORM
- [PostgreSQL](https://www.postgresql.org/) - Database
- [Redis](https://redis.io/) - Task Queue, Cache
- [Asynqmon](https://github.com/hibiken/asynqmon) - Task Queue Monitoring

## Getting Started

### Clone the repository
```shell
git clone https://github.com/vuquang23/poseidon.git
cd poseidon
```

### Start infra
```shell
docker-compose -f ./infra/docker-compose-infra.yaml -p poseidon-infra up -d
```

### Run from source
Make sure you have Go installed ([download](https://go.dev/dl/)). Encourage using version `1.21` or higher. <br/>
**NOTE**: `API Server` needs to be run first for database migration.

#### Install dependencies
```shell
go mod tidy
```

#### Run API Server

```shell
go run cmd/app/main.go api
```

#### Run Master
```shell
go run cmd/app/main.go master
```

#### Run Worker
```shell
go run cmd/app/main.go worker
```

### Run in containers
Make sure you've already started the infra.

#### Build image locally
```shell
docker build . -t poseidon
```

#### Start all components
Mac/Window

```shell
docker-compose -p poseidon-app up -d
```

Linux

```shell
docker-compose -f docker-compose-linux.yaml -p poseidon-app up -d
```

### Testing
Need Go installed ([download](https://go.dev/dl/)) with version `1.21` or higher. <br/>

#### Start infra for testing
```shell
docker-compose -f ./infra/docker-compose-infra-test.yaml -p poseidon-infra-test up -d
```

#### Run tests
```shell
go test -p 1 -coverpkg=./... -coverprofile=profile.cov ./...
```

#### Test coverage
```shell
go tool cover -func profile.cov

go tool cover -html profile.cov
```

## Documentation
Please checkout [Architecture](docs/architecture.md) for more details on system architecture, and [API Docs](docs/api.yaml) for API documentation.
