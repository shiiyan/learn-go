# docker runtime for learning go

## how to use

pull the latest image.

```
docker compose pull
```

start a new container and opening a shell inside it.

```
docker compose run --remove-orphans app bash
```

move to target package.

```
cd <target-package>
```

run go file.

```
go run main.go
```
