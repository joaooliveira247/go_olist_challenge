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

<details>
    <summary><code>/authors/</code></summary>

<code>POST /authors/</code>

- **Description**: Creates a new author.

- **Headers**:

    ```plaintext
    Content-Type: application/json

    ```

- **Request Body**:

    ```json
    {
        "name": "Stephen King"
    }
    ```

- **Success Response (201 Created)**:

    ```json
    {
         "id": "1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47"
    }
    ```

- **Errors**:

    - **422 Unprocessable Entity**: Invalid request body.

    - **409 Conflict**: Author already exists.

    - **500 Internal Server Error**: Failed to create the entity.

- **Example Request with cURL**:

```curl
curl -X POST localhost:8000/authors/ \
-H "Content-Type: application/json" \
-d '{
    "name": "Stephen King"
}'
```

<code>GET /authors/</code>

- **Description**: Retrieves authors based on query parameters. If no parameters are provided, it returns all authors.

- **Headers**:

    ```plaintext    
    Content-Type: application/json
    ```

- **Query Parameters**:

    **authorID** (string, optional): UUID of the author.

    **name** (string, optional): Name of the author.

- **Success Responses (200 OK)**:

    - Single Author by ID.

        ```json
        {
            "id":   "1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47",
            "name": "Stephen King"
        }
        ```

    - Multiple Authors by Name

        ```json
        [
            {
                "id": "1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47",
                "name": "Stephen King"
            },
            {
                "id": "2a8c2dde-24b3-4c21-9fbb-d7dfd09f98e5",
                "name": "Stephen Hawking"
            }
        ]
        ```

    - All Authors

        ```json
        [
            {
                "id": "1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47",
                "name": "Stephen King"
            },
            {
                "id": "2a8c2dde-24b3-4c21-9fbb-d7dfd09f98e5",
                "name": "J.K. Rowling"
            }
        ]
        ```

- **Errors**:

    - **400 Bad Request**: Invalid query parameters or invalid ID.

    - **404 Not Found**: Author not found.

    - **500 Internal Server Error**: Unable to fetch entity.

- **Example Requests with cURL**:

    - Get All Authors

        ```bash
        curl -X GET localhost:8000/authors/ \
        -H "Content-Type: application/json"
        ```

    - Get Author by ID

        ```bash
        curl -X GET "localhost:8000/authors/?authorID=1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47" \
        -H "Content-Type: application/json"
        ```

    - Get Authors by Name

        ```bash
        curl -X GET "localhost:8000/authors/?name=Stephen" \
        -H "Content-Type: application/json"
        ```

<code>DELETE /authors/{id}</code>

- **Description**: Deletes an author by ID.

- **Headers**:

    ```plaintext
    Content-Type: application/json
    ```

- **Path Parameter**:

**id** (string, required): UUID of the author to be deleted.

- **Success Response (204 No Content)**:

```json
(empty response body)
```

- **Errors**:

    - **400 Bad Request**: Invalid author ID.

    - **404 Not Found**: Author not found.

    - **500 Internal Server Error**: Unable to delete the entity.

- **Example Request with cURL**:

```bash
curl -X DELETE localhost:8000/authors/1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47 \
  -H "Content-Type: application/json"
```

</details>

<details>
<summary><code>/books/</code></summary>
</details>

## 📦 Usage libraries:

