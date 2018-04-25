Log event with GRPC
- Send event to server
- Query events from server

# Config
Edit file `config.yaml`
# Init db
Run following command
```sh
$ go run migrate/*.go
```
# Build
```sh
$ make
```

# Gen GRPC
```sh
# make grpc
```
# Usage
1. Build
```sh
$ make
```
2. Edit config file `config.yaml`
3. Start grpc server
```sh
$ ./bin/server/elog
```
4. Create/query
```sh
$ ./bin/client/elog -h // for help
```