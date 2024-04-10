# A project that I have written in Go

```shell
go list -m all
```

This project use **swagger**

Run this command to install `swag` command

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

Init swag

```shell
swag init
```

Download lib

```shell
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

## Upgrade project lib

To update go project

```shell
go get -u
go mod tidy
```

One can run the command to view available upgrades for direct dependencies. Unfortunately, the output is not actionable, i.e. we can't easily use it to update multiple dependencies.

```shell
go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null
```
