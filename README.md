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

Install the binary as a global command.

```bash
go install github.com/aminnairi/gouache
```

Run your request.

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

> [!WARNING]
> You should have the path to go installed binaries [already
> setup](https://go.dev/wiki/SettingGOPATH) in order for your terminal to have
> access to this program.

## 👋 Uninstallation

> [!NOTE]
> The `go` binary does not have a way to uninstall a previously installed
> package, but all it does is create a folder and download the sources in that
> folder.

```bash
rm -rf $(which gouache)
```

## ❓ Documentation

### Request

#### Provide an HTTP request from a file

> [!IMPORTANT]
> If you are unsure of what the HTTP protocol is, [here is a detailed
> article](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Overview)
> from the Mozilla Developers Network website that explains it well.

```bash
touch get.http
```

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com/users
```

> [!WARNING]
> There should only be one HTTP request per `*.http` file.

```bash
gouache request get.http
```

#### Display the status line of the response

```bash
gouache request --with-status get.http
# or
gouache request -s get.http
```

#### Display the headers of the response

```bash
gouache request --with-headers get.http
# or
gouache request -H get.http
```

#### Display the body of the request

```bash
gouache request --with-body get.http
# or
gouache request -b get.http
```

#### Send a POST request

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
gouache request post.http
```

#### Send a PATCH request

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
gouache request patch.http
```

#### Send a PUT request

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

#### Send a DELETE request

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
gouache request delete.http
```

#### Run all requests from a folder recursively

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

### Generate

#### Generate a request

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

> [!TIP]
> Only providing one of the following: method, path or host will trigger the
> interactive mode asking you for the remaining informations.

#### Generate a request with a body

```bash
gouache generate post.http \
  --method GET \
  --path /users \
  --host https://jsonplaceholder.typicode.com \
  --body '{"id":1}'
```

> [!TIP]
> Requests don't necessarily have to have a body, omitting it will not trigger
> the interactive mode and rather will create the request without body.

#### Generate a request in interactive mode

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
