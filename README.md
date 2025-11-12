# gouache

Create HTTP requests from file and run them right from your terminal

## Synopsis

`gouache`, pronounced _/ɡwɑːʃ/_, it represent a type of paint that can react to
water to update or create a new painting.

## Usage

```bash
touch index.request
```

```http
GET /users HTTP/2
Host: https://jsonplaceholder.typicode.com
Accept: application/json
```

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

```bash
go install github.com/aminnairi/gouache
```

## License

See [`LICENSE`](./LICENSE).
