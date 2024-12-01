Wallet Service
==

> Author: Liu <lightserver.cn@gmail.com>
>
> Date: 2024-12-01


### Project Overview

This project aims to implement a wallet service that meets specific user stories and provides services through a RESTful API. The project is developed using Go language and PostgreSQL database.

### 环境准备

- Go（Version 1.23.3）
- PostgreSQL
- Docker
- Docker Compose

### Code Structure

The code organization structure is as follows:

```sh
- app: Contains business logic and handlers, which is the core part of the project
  - controller: Handles HTTP requests and responses, distributes requests, and returns service responses
  - model: Defines data structures and database models
  - repository: Handles interactions with the database, provides data operation interfaces
  - request: Defines the structure of API requests
  - service: Implements business logic and rules, calls repository for data operations
- boot: Initializes related components
  - boot: Initialization startup code
  - config: Loads and parses configuration files
  - db: Initializes database connections
  - http: Initializes HTTP server
  - log: Initializes logging system
- cmd: Contains the main application entry point
- config: Contains configuration files
  - config.go: Configuration file structure definition
  - config.local.yaml: Local configuration file
  - config.yaml: Container configuration file
  - ddl.sql: Database table structure DDL
- docker-compose: Defines services and their dependencies
  - volumes: Data volumes
    - postgres: PostgreSQL data volume
    - redis: Redis data volume
- dockerfile: Builds Docker images
- pkg: Contains reusable packages and utilities
  - consts: Constant definitions
  - dal: Data access layer
  - log: Custom logging package
- postman: API testing interface collection
- router: Defines API routes
- runtime: Runtime configuration and scripts
- test: Unit tests and integration tests
  - db: Database-related tests
```

### Installation and Running

#### Local Running

1. Set environment variables:

Edit the `config/config.local.yaml` file to set necessary environment variables.

Example:

```shell
db:
  driver: postgres
  host: 127.0.0.1
  port: 5432
  user: postgres
  password: postgres
  db_name: postgres

db-test:
  driver: postgres
  host: 127.0.0.1
  port: 5432
  user: postgres
  password: postgres
  db_name: test_postgres

redis:
  addr: 127.0.0.1:6379
  password:
  db: 0
```

2. Run the application:

```shell
go run main.go
```

#### Containerized Running

Execute the following command:

```shell
sh ./run.sh
```

The command includes:
- Building service binaries
- Creating `runtime/log` log directory
- Creating `docker-compose/volumes/{postgres,redis}`
- Building a runtime environment containing services (postgres, adminer, redis, redis-commander, golang)

#### API Testing

1. Access `http://127.0.0.1:8082/` and log in to pgsql

> Enter the account and password according to the docker-compose/docker-compose.yaml configuration

![PostgreSQL](./pics/PostgreSQL.png)

2. Select `SQL command` to execute and initialize the project data table structure

> If db.init_table is configured to be enabled by default in config/config.yaml, there is no need to build the data table structure.

![SqlCommand](./pics/SqlCommand.png)

3. Open Postman, import data, and conduct testing

> Data is in postman/wallet.postman_collection.json

![Postman](./pics/Postman.png)

4. Register two initial users for testing purposes, send a POST request to http://localhost:8080/api/users.

5. Wallet-related interfaces are in the wallets folder.

### Decision Description

- Language: Go is chosen for its performance, concurrency features, and powerful standard library.
- Database: PostgreSQL is chosen for its robustness, support for ACID transactions, and rich feature set.
- In-memory database: Redis is used for caching and handling high-frequency data access.
- ORM: Not allowed according to requirements, so raw SQL queries with database/sql package are used.
- Logging: go.uber.org/zap package is used to simplify operations.
- Decimal handling: github.com/shopspring/decimal package is used for precise decimal calculations.


### Linting

The project uses golangci-lint for linting. The root directory provides a golangci.yaml configuration file. Run linting:

```shell
golangci-lint run --config=.golangci.yaml
```

### Unit Testing

Unit tests are written using Go's built-in testing package. Run tests:

In some cases, the following error may occur, which is not resolved yet but does not affect business testing

> CreateTestDB returned an error: pq: duplicate key value violates unique constraint "pg_database_datname_index"

```shell
go test ./... -race -cover
```

### Goroutine Leak Checking

Use the uber-go/goleak package to check for Goroutine leaks. Run leak checking:

In some cases, the following error may occur, which is not resolved yet but does not affect business testing

> CreateTestDB returned an error: pq: duplicate key value violates unique constraint "pg_database_datname_index"

```shell
go test ./... -race -cover -timeout=30m -v -tags=leakcheck
```

### Time Estimation

Below is the estimated time that spent on developing and maintaining this project:

- Requirements analysis and design: 4 hours
- Environment setup and configuration: 4 hours
- Project setup and running: 12 hours
- Core feature implementation (including business logic, database interaction, API development): 20 hours
- Unit test and integration test writing: 32 hours
- Dockerization and container orchestration: 8 hours
- Documentation writing and README optimization: 4 hours
- Total: 84 hours
