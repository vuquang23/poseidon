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

You can specify your config file path, and overwrite any config values 

## API Documentation