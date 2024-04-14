# Usage

## CLI

```
NAME:
   Poseidon - A new cli application

USAGE:
   Poseidon [global options] command [command options]

COMMANDS:
   api      Run API Server
   worker   Run Worker
   master   Run Master
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  Configuration file (default: "internal/pkg/config/default.yaml")
   --help, -h                show help
```

**NOTE**: check CLI usage by the log after building Docker image, or ```go run cmd/app/main.go```.


## Configuration

Default config is stored at [default.yaml](/internal/pkg/config/default.yaml).

You can specify your config file path, and override any config values in the file by env var.


For examples:
```shell
export ETH_RPC=https://example-archieve-node.com

export POSTGRES_HOST=postgresql.db-cluster.gke.develop.internal

export SERVICE_TASK_BLOCKFINALIZATION=100
```

## API Documentation

[openapi.yaml](/docs/openapi.yaml) hosted on [stoplight.io](https://psdon.stoplight.io/docs/poseidon/16e3c17d0c7f5-poseidon)
