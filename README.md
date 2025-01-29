# [📖 Work at Olist Challenge](https://github.com/olist/work-at-olist)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/joaooliveira247/go_olist_challenge)
![GitHub Workflow Tests](https://github.com/joaooliveira247/go_olist_challenge/actions/workflows/run-tests.yaml/badge.svg)

## 💻 Requirements:

### `Go >= 1.22.5`

### [`Docker`](https://www.docker.com/) & [`Docker compose`](https://docs.docker.com/compose/)

- Installation & Running:

    <details>

    <summary><code>Docker</code></summary>

    - Starts all services, including the API and Database.

        ```bash
        docker compose up -d
        ```

    - Create all tables in database.

        ```bash
        make db create CONTAINER_ID=<container_id>
        ```

    - Delete all tables in database.

        ```bash
        make db delete CONTAINER_ID=<container_id>
        ```

    - Imports an authors CSV file into the database. Ensure the CSV file includes headers.

        ```bash
        make import CSV_PATH=<csv_path> CONTAINER_ID=<container_id>
        ```

        > **NOTE:**
        >
        > To find the container_id, run `docker ps`
    
    </details>

    <details>
    <summary><code>Local</code></summary>

    - Install all packages.

        ```bash
        go mod tidy
        ```

    - Run API.

        ```bash
        air run
        ```

    - Create all tables.

        ```bash
        go run main.go db create
        ```

    - Delete all tables.

        ```bash
        go run main.go db delete
        ```

    - Imports an authors CSV file into the database.

        ```bash
        go run main.go <path_csv> --header <true|false>
        ```

    </details>

## 📜 Documentation:

## 📦 Usage libraries:

