# gouache

Create HTTP requests from file and run them right from your terminal

> [!NOTE] > `gouache`, pronounced _/ɡwɑːʃ/_, it represent a type of paint that can react to water to update or create a new painting.

## Usage

```bash
touch index.request
```

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com
Accept: application/json
```

> [!WARNING]
> There should only be one HTTP request per `*.request` file.

```bash
gouache -request index.request -with-status -with-headers -with-body
HTTP/2 200 OK
Server: cloudflare
X-Powered-By: Express
Content-Type: application/json; charset=utf-8
Cache-Control: max-age=43200
[
  {
    "id": 1,
    "name": "Leanne Graham",
    "username": "Bret",
    "email": "Sincere@april.biz",
    "address": {
      "street": "Kulas Light",
      "suite": "Apt. 556",
      "city": "Gwenborough",
      "zipcode": "92998-3874",
      "geo": {
        "lat": "-37.3159",
        "lng": "81.1496"
      }
    },
...
```

## Installation

> [!NOTE]
> This commands needs to be run inside of a terminal with access to the `go` binary.

```bash
go install github.com/aminnairi/gouache
```

## FAQ

### How to create a GET request

```bash
touch get.request
```

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com/users
```

```bash
gouache -request get.request -with-status
```

### How to create a POST request

```bash
touch post.request
```

```http
POST /users HTTP/2
Host: https://jsonplaceholder.typicode.com/users
Content-Type: application/json

{
  "id": 11,
  "email": "user@domain.com"
}
```

```bash
gouache -request post.request -with-status
```

### How to create a PATCH request

```bash
touch patch.request
```

```http
PATCH /users/10 HTTP/2
Host: https://jsonplaceholder.typicode.com/users
Content-Type: application/json

{
  "email": "user@domain.com"
}
```

```bash
gouache -request patch.request -with-status
```

### How to create a PUT request

```bash
touch put.request
```

```http
PUT /users/10 HTTP/2
Host: https://jsonplaceholder.typicode.com/users
Content-Type: application/json

{
  "email": "user@domain.com"
}
```

```bash
gouache -request put.request -with-status
```

### How to create a DELETE request

```bash
touch delete.request
```

```http
DELETE /users HTTP/2
Host: https://jsonplaceholder.typicode.com/users
Content-Type: application/json

{
  "id": 10,
  "email": "user@domain.com"
}
```

```bash
gouache -request delete.request -with-status
```

### Run a request without output

```bash
touch create-user.request
```

```http
POST /users HTTP/2
Host: https://jsonplaceholder.typicode.com/users
Content-Type: application/json

{
  "id": 11,
  "email": "user@domain.com"
}
```

```bash
gouache -request create-user.request
```

> [!IMPORTANT]
> Errors will still be written in the standard error of your terminal.

## License

See [`LICENSE`](./LICENSE).
