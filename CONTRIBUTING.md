# Contributing

Contribution guidelines for developers

> [!TIP]
> Looking for the documentation on how to run this project as a user? Head over [here](./README.md).

## Requirements

- [`git`](https://git-scm.com/)
- [`go`](https://go.dev/dl/)

## Clone

> [!IMPORTANT]
> If you want to participate in the development and submit pull requests, you should rather [fork](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo) this project and use your own copy.

```bash
git clone https://github.com/aminnairi/gouache
cd gouache
```

## Run the program

```bash
go run . --with-status examples
```

## Build the program

```bash
go build
./gouache --with-status examples
```

## Install the program as a global binary

> [!WARNING]
> You should [add go in your environment variable path](https://go.dev/wiki/GOPATH) in order to have access to it as a global command in your shell.

```bash
go install
```

## Uninstall the globally installed binary

```bash
rm -rf $(which gouache)
```
