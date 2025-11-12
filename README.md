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
GET /1.1.1.1/json HTTP/2
Host: ipapi.co
Accept: application/json
```

```bash
gouache --request index.http
HTTP/2 200 OK
```

## Installation

```bash
go install github.com/aminnairi/gouache
```

## License

See [`LICENSE`](./LICENSE).
