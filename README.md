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


<code>/authors/</code>

<details>
<summary><code>POST /authors/</code></summary>

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
</details>

<details>
<summary><code>GET /authors/</code></summary>

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

</details>

<details>
<summary><code>DELETE /authors/{id}</code></summary>

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

###

<code>/books/</code>

<details>
<summary><code>POST /books/</code></summary>

- **Description**: Creates a new book with its associated authors.

- **Headers**:

    ```plaintext
    Content-Type: application/json
    ```

- **Request Body**:

    ```json
    {
        "title": "The Shining",
        "edition": 1,
        "published_year": 1977,
        "authors_id": [
            "1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47"
        ]
    }
    ```

- **Success Response (201 Created)**:

    ```json
    {
        "id": "3f8c3bde-54a6-41d7-bb4f-8d74a33e8e12"
    }
    ```

- **Errors**:

    - **422 Unprocessable Entity**: Invalid request body.

    - **500 Internal Server Error**: Failed to create the entity.

- **Example Request with cURL**:

    ```bash
    curl -X POST localhost:8000/books/ \
    -H "Content-Type: application/json" \
    -d '{
        "title": "The Shining",
        "edition": 1,
        "published_year": 1977,
        "authors_id": [
        "1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47"
        ]
    }'
    ```

</details>

<details>
<summary><code>GET /books/</code></summary>

- **Description**: Retrieves a list of books. Supports filtering by book ID, author ID, title, edition, and publication year.

- **Headers**:

    ```plaintext
    Content-Type: application/json
    ```

- **Query Parameters**:

    **bookID** (optional, UUID): Filters books by their unique ID.

    **authorID** (optional, UUID): Filters books by the author's unique ID.

    **title** (optional, string): Filters books by title (case-insensitive).

    **edition** (optional, uint8): Filters books by edition number.

    **publicationYear** (optional, uint): Filters books by publication year.

    **title, edition, and publicationYear** can be used together for a more precise query.

- **Success Response (200 OK)**:

    ```json
    [
        {
            "id": "3f8c3bde-54a6-41d7-bb4f-8d74a33e8e12",
            "title": "The Shining",
            "edition": 1,
            "publicationYear": 1977,
            "authors": [
                    "name": "Stephen King"
            ]
        }
    ]
    ```

- **Errors**:

    - **400 Bad Request**: Invalid query parameter.

    - **400 Bad Request**: Invalid ID format.

    - **500 Internal Server Error**: Failed to fetch the entity.

- **Example Requests with cURL**:

    - **Get all books**:
        ```bash
        curl -X GET "localhost:8000/books/" -H "Content-Type: application/json"
        ```

    - **Get a book by ID**:
        ```bash
        curl -X GET "localhost:8000/books/?bookID=3f8c3bde-54a6-41d7-bb4f-8d74a33e8e12" \
        -H "Content-Type: application/json"
        ```

    - **Get books by author ID**:

        ```bash
        curl -X GET "localhost:8000/books/?authorID=1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47" \
        -H "Content-Type: application/json"
        ```

    - **Get books by title**:

        ```bash
        curl -X GET "localhost:8000/books/?title=The%20Shining" \
        -H "Content-Type: application/json"
        ```

    - **Get books by title, edition, and publication year**:

        ```bash
        curl -X GET "localhost:8000/books/?title=The%20Shining&edition=1&publicationYear=1977" \
        -H "Content-Type: application/json"
        ```

</details>

<details>
<summary><code>PUT /books/{id}</code></summary>

- **Description**: Updates an existing book by its ID. Allows modifying book details and authors.

- **Headers**:

    ```plaintext
    Content-Type: application/json
    ```

- **Path Parameter**:

    **id** (UUID, required): The unique identifier of the book to update.

- **Request Body** (partial or full update allowed):

    ```json
    {
        "title": "Updated Title",
        "edition": 2,
        "publicationYear": 2024,
        "authorsID": [
            "1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47",
            "9a6c112e-fc2e-49d3-b930-7991a20903db"
        ]
    }
    ```

- **Success Response (204 No Content)**:

```json
(empty response body)
```

- **Errors**:

    - **400 Bad Request**: Invalid ID format.

    - **422 Unprocessable Entity**: Invalid request body.

    - **304 Not Modified**: Nothing to update.

    - **500 Internal Server Error**: Failed to fetch or update the entity.

- **Example Requests with cURL**:

    - **Update book details**:

        ```bash
        curl -X PUT "localhost:8000/books/3f8c3bde-54a6-41d7-bb4f-8d74a33e8e12" \
        -H "Content-Type: application/json" \
        -d '{
            "title": "Updated Title",
            "edition": 2,
            "publicationYear": 2024
        }'
        ```

    - **Update book authors**:

        ```bash
        curl -X PUT "localhost:8000/books/3f8c3bde-54a6-41d7-bb4f-8d74a33e8e12" \
        -H "Content-Type: application/json" \
        -d '{
            "authorsID": [
            "1d47bbe5-c7d3-4580-ad2a-c4b192eeeb47",
            "9a6c112e-fc2e-49d3-b930-7991a20903db"
            ]
        }'
        ```

</details>

<details>
<summary><code>DELETE /books/{id}</code></summary>

- **Description**: Deletes a book by its ID.

- **Path Parameter**:

    **id** (UUID, required): The unique identifier of the book to delete.

- **Success Response (204 No Content)**:

```json
(empty response body)
```
- Errors:

    - **400 Bad Request**: Invalid ID format.

    - **304 Not Modified**: Book not found (nothing to delete).

    - **500 Internal Server Error**: Unable to fetch or delete the entity.

- **Example Request with cURL**:

```bash
curl -X DELETE "localhost:8000/books/3f8c3bde-54a6-41d7-bb4f-8d74a33e8e12"
```

</details>

## 📦 Usage libraries:

