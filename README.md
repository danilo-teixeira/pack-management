# pack-management

This project uses asdf to manage the dependences: https://github.com/asdf-vm/asdf

Install deps
```sh
asdf install && asdf reshim golang
````

Run local database image
```sh
docker compose -f './docker-compose.local.yml' up -d --build 'db'
```
