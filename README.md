# [📖 Work at Olist Challenge](https://github.com/olist/work-at-olist)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/joaooliveira247/go_olist_challenge)
![GitHub Workflow Tests](https://github.com/joaooliveira247/go_olist_challenge/actions/workflows/run-tests.yaml/badge.svg)

## 💻 Requirements:

### `Go >= 1.22.5`

### [`Docker`](https://www.docker.com/) & [`Docker compose`](https://docs.docker.com/compose/)

- Installation & Running:

    <details>

    <summary><code>Docker</code></summary>

        `docker compose up -d`

            Starts all services, including the API and Database.

        `make db create CONTAINER_ID=<container_id>`

            Create all tables in database.

        `make db delete CONTAINER_ID=<container_id>`

            Delete all tables in database.

        `make import CSV_PATH=<csv_path> CONTAINER_ID=<container_id>`

              Imports an authors CSV file into the database. Ensure the CSV file includes headers.

        > **NOTE:**
        >
        > To find the container_id, run `docker ps`
    
    </details>

## 📜 Documentation:

## 📦 Usage libraries:

