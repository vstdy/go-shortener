# Shortener

The application is a tool that creates a short, unique URL that will redirect to the specific website of your choosing.  
Is a project of the [course](https://practicum.yandex.ru/promo/go-profession/)

## REST API

By default, server starts at `8080` HTTP port with the following endpoints:

- `POST /` - create shortcut for url from text/plain body;
- `POST /api/shorten` - create shortcut for url from application/json body;
- `POST /api/shorten/batch` - create shortcuts for urls batch from application/json body;
- `GET /{urlID}` - follow origin url from shortcut;
- `GET /api/user/urls` - get urls created by current user;
- `DELETE /api/user/urls` - remove urls created by current user with given ids;
- `GET /ping` - check connection to database;

For details check out [***http-client.http***](./http-client.http) file


## Code
### Packages used

App architecture and configuration:

- [viper](https://github.com/spf13/viper) - app configuration;
- [cobra](https://github.com/spf13/cobra) - CLI;
- [zerolog](https://github.com/rs/zerolog) - logger;
- [jaeger-client-go](https://github.com/jaegertracing/jaeger-client-go) - tracer;

Networking:

- [go-chi](https://github.com/go-chi/chi) - HTTP router;
- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) - gRPC to JSON proxy generator;

SQL database interface provider:

- [bun](https://github.com/uptrace/bun) - database client;

Testing:

- [testify](https://github.com/stretchr/testify) - tests build toolkit;
- [testcontainers-go](https://github.com/testcontainers/testcontainers-go) - library to create runtime environment for tests;

## CLI

All CLI commands have the following flags:
- `--log_level`: (optional) logging level (default: `info`);
- `--config`: (optional) path to configuration file (default: `./config.toml`);
- `--timeout`: (optional) request timeout (default: `5s`);
- `-d --database_dsn`: (optional) database source name (default: `postgres://user:password@localhost:5432/shortener?sslmode=disable`);

Root only command flags:
- `-a --server_address`: (optional) server address (default: `0.0.0.0:8080`);
- `-b --base_url`: (optional) base URL (default: `http://127.0.0.1:8080`);
- `-s --storage_type`: (optional) storage type (default: `psql`);
- `-f --file_storage_path`: (optional) file storage path (default: `./storage/file/storage_file.txt`);

If config file not specified, defaults are used. Defaults can be overwritten using ENV variables.

### Migrations

    shortener migrate --config ./my-confs/config-1.toml

Command migrates DB to the latest version

### gRPC client
1. Shorten given url.
   ```
   shortener client shorten https://lengthy-url.com/
   ```
   Arguments:
   - `args[0]`: url to shorten;

1. Shorten given urls batch.
    ```
    shortener client shorten batch https://lengthy-url-1.com/ https://lengthy-url-2.com/
    ```
   Arguments:
   - `args[0] args[1]...`: urls to shorten

1. Return original url from shortened url.
    ```
    shortener client get original_url http://127.0.0.1:8080/gw/1
    ```
   Arguments:
   - `args[0]`: shortened url;

1. Return user's urls.
    ```
    shortener client get users_urls
    ```
1. Delete user's urls.
    ```
    shortener client delete 1 3
    ```
   Arguments:
   - `args[0] args[1]...`: ids of urls to delete

Flags:
- `-t --token`: (optional) user's token;

## How to run
### Docker

    docker-compose -f build/docker-compose.yml up
