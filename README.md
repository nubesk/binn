# binn
binn is open-source message in a bottle server.

[![Test](https://github.com/nubesk/binn/actions/workflows/test.yml/badge.svg)](https://github.com/nubesk/binn/actions/workflows/test.yml)

## usage
### build image
```
docker build -t binn .
```

### run container
```
docker run -it --rm -v $(pwd):/go/src/github.com/binn -p 8080:8080 binn /bin/sh
```

### run app
```
go run app.go
```
