# Shortener

The application is a tool that creates a short, unique URL that will redirect to the specific website of your choosing.  
Is a project of the [course](https://practicum.yandex.ru/profile/go-developer/)

## REST API

By default, server starts at `8080` HTTP port with the following endpoints:

- `POST /` - create shortcut for url from text/plain body;
- `POST /api/shorten` - create shortcut for url from application/json body;
- `POST /api/shorten/batch` - create shortcuts for urls batch from application/json body;
- `GET /{urlID}` - follow origin url from shortcut;
- `GET /api/user/urls` - get urls created by current user;
- `DELETE /api/user/urls` - remove urls created by current user with given ids;
- `GET /ping` - check connection to database;

For details check out [***requests.http***](./requests.http) file


## Code
### Packages used

App architecture and configuration:

- [viper](https://github.com/spf13/viper) - app configuration;
- [cobra](https://github.com/spf13/cobra) - CLI;
- [zerolog](github.com/rs/zerolog) - logger;

Networking:

- [go-chi](github.com/go-chi/chi) - HTTP router;

SQL database interface provider:

- [bun](github.com/uptrace/bun) - database client;

Testing:

- [testify](github.com/stretchr/testify) - tests build toolkit;
- [testcontainers-go](github.com/testcontainers/testcontainers-go) - library to create runtime environment for tests;

## CLI

All CLI commands have the following flags:
- `--log-level`: (optional) set logging level (default: `info`);
- `--config`: (optional) path to configuration file (default: `./config.toml`);
- `--timeout`: (optional) request timeout (default: `5s`);
- `-d --database_dsn`: (optional) database source name (default: `postgres://user:password@localhost:5432/shortener?sslmode=disable`);

Root only command flags:
- `-a --server_address`: (optional) server address (default: `0.0.0.0:8080`);
- `-b --base_url`: (optional) base URL (default: `http://127.0.0.1:8080`);
- `-s --storage_type`: (optional) set logging level (default: `psql`);
- `-f --file_storage_path`: (optional) set logging level (default: `./storage/file/storage_file.txt`);

### Migrations

    shortener migrate --config ./my-confs/config-1.toml

If config file not specified, defaults are used. Defaults can be overwritten using ENV variables.

## How to run
### Docker

    docker-compose -f build/docker-compose.yml up
