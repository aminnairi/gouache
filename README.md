# 🖌️ gouache

Create HTTP requests from files and run them right from your terminal

> [!NOTE]
> The name of this library, `gouache`, pronounced _/ɡwɑːʃ/_, is a french word representing a type
> of paint that can react to water to update or create a new painting.

## 🤔 Usage

Create a new HTTP request.

```bash
touch index.http
```

Fill in the details of your request.

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com
Accept: application/json
```

> [!WARNING]
> There should only be one HTTP request per `*.http` file.

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
gouache request --with-status --with-headers --with-body index.http
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

## ✨ Features

- Write HTTP requests in a friendly format
- It's just HTTP protocol and nothing else
- Run requests right from your terminal
- Run all requests from a folder recursively

## 🏃 Installation

> [!NOTE]
> This commands needs to be run inside of a terminal with access to the `go` binary.

```bash
go install github.com/aminnairi/gouache
```

## 👋 Uninstallation

> [!NOTE]
> The `go` binary does not have a way to uninstall a previously installed
> package, but all it does is create a folder and download the sources in that
> folder.

```bash
rm -rf $(which gouache)
```

## ❓ FAQ

> [!IMPORTANT]
> If you are unsure of what the HTTP protocol is, [here is a detailed
> article](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Overview)
> from the Mozilla Developers Network website that explains it well.

### How to create a GET request

```bash
touch get.http
```

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com/users
```

```bash
gouache request --with-status get.http
```

### How to create a POST request

```bash
touch post.http
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
gouache request --with-status post.http
```

### How to create a PATCH request

```bash
touch patch.http
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
gouache request --with-status patch.http
```

### How to create a PUT request

```bash
touch put.http
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
gouache request --with-status put.http
```

### How to create a DELETE request

```bash
touch delete.http
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
gouache request --with-status delete.http
```

### Run a request without output

```bash
touch create-user.http
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
gouache request create-user.http
```

> [!IMPORTANT]
> Errors and informations will still be written in the standard error of your terminal.

### Run all requests from a folder recursively

```bash
mkdir requests
touch requests/users.http
```

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com
```

```bash
touch requests/posts.http
```

```http
GET /posts HTTP/2
Host: https://jsonplaceholder.typicode.com
```

```bash
mkdir requests/posts
touch requests/posts/first.http
```

```http
GET /posts/1 HTTP/2
Host: https://jsonplaceholder.typicode.com
```

```bash
gouache request --with-status requests
```

### Generate a request right from the command line

```bash
gouache generate get.http \
  --method GET \
  --path /users \
  --host https://jsonplaceholder.typicode.com
# or
gouache generate get.http \
  -m GET \
  -p /users \
  -H https://jsonplaceholder.typicode.com
```

### Generate a request with a body

```bash
gouache generate post.http \
  --method GET \
  --path /users \
  --host https://jsonplaceholder.typicode.com \
  --body '{"id":1}'
```

### Generate a request in interactive mode

```bash
gouache generate
```

## ⚖️ License

See [`LICENSE`](./LICENSE).

## 🤝 Contributing

See [`CONTRIBUTING.md`](./CONTRIBUTING.md).

## 🔒 Security

See [`SECURITY.md`](./SECURITY.md).

## 🫶 Code of conduct

See [`CODE_OF_CONDUCT.md`](./CODE_OF_CONDUCT.md).
