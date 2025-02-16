# pack-management

This project uses asdf to manage the dependences: https://github.com/asdf-vm/asdf

Install Golang using asdf:

```sh
asdf install && asdf reshim golang
```

## Run local

1. Copy .env files:

```sh
make init
```

2. Set environment variables in the .env and .env.test files;
3. Install deps:

```sh
make install
```

4. Run:

```sh
make run
```

## Run integration tests

```sh
make test-e2e
```

## Run load test
To run load teste go to [load test folder](./__loadtest/README.md)
