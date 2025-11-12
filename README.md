# gouache

Create HTTP requests from file and run them right from your terminal

> [!NOTE]
> The name of this library, `gouache`, pronounced _/ɡwɑːʃ/_, is a french word representing a type
> of paint that can react to water to update or create a new painting.

## Usage

Create a new HTTP request.

```bash
touch index.request
```

Fill in the details of your request.

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com
Accept: application/json
```

> [!WARNING]
> There should only be one HTTP request per `*.request` file.

Install the binary as a global command.

```bash
go install github.com/aminnairi/gouache
```

Run your request.

> [!WARNING]
> You should have the path to go installed binaries [already
> setup](https://go.dev/wiki/SettingGOPATH) in order for your terminal to have
> access to this program.

```bash
gouache --request index.request --with-status --with-headers --with-body
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

## Features

- Write HTTP requests in a friendly format
- It's just HTTP protocol and nothing else
- Run requests right from your terminal

## Installation

> [!NOTE]
> This commands needs to be run inside of a terminal with access to the `go` binary.

```bash
go install github.com/aminnairi/gouache
```

## Uninstallation

> [!NOTE]
> The `go` binary does not have a way to uninstall a previously installed
> package, but all it does is create a folder and download the sources in that
> folder.

```bash
rm -rf $(which gouache)
```

## FAQ

> [!IMPORTANT]
> If you are unsure of what the HTTP protocol is, [here is a detailed
> article](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Overview)
> from the Mozilla Developers Network website that explains it well.

### How to create a GET request

```bash
touch get.request
```

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com/users
```

```bash
gouache --request get.request --with-status
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
gouache --request post.request --with-status
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
gouache --request patch.request --with-status
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
gouache --request put.request --with-status
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
gouache --request delete.request --with-status
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
gouache --request create-user.request
```

> [!IMPORTANT]
> Errors will still be written in the standard error of your terminal.

## License

See [`LICENSE`](./LICENSE).
