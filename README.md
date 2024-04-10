# Poseidon

## Setting

### Dependencies

```
$ go mod tidy
```

### Infra

```
$ docker-compose -f ./infra/docker-compose-infra.yaml -p poseidon-infra up -d
```

## Run API Server

```
$ go run cmd/app/main.go api
```

## Run Taskq Master

```
$ go run cmd/app/main.go master
```

## Run Taskq Worker

```
$ go run cmd/app/main.go worker
```

## Testing

```
$ docker-compose -f ./infra/docker-compose-test.yaml -p poseidon-infra-test up -d
$ make test
$ docker-compose -f ./infra/docker-compose-test.yaml -p poseidon-infra-test down
```
