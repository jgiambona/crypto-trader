# CryptoTrader

CryptoTrader automates the task of trading cryptocurrencies based on current market price.

***

### Current State

+ LiveCoin Exchange API
+ Web frontend (WIP)

### What's in the Roadmap?

- Enable and disable crypto exchange.
- Control the platform using telegram.
- Need more crypto exchange.

### Production

#### Prerequisites:

+ Docker 18.1

The first thing to do in order to run system is to download docker and pull golang image into the host system.

~~~
docker pull golang:1.10.0
docker run --rm -it -p 4000:4000 -v $PWD:/app golang:1.10.0 bash -C "go run main.go init_*.go"
~~~

### Development

#### Prerequisites:

+ Go 1.10

Assuming, Go is already in the system execute this in the workspace root.

~~~
go run main.go init_*.go
~~~

Set the configuration `config.yaml` and restart.


~~~
$ docker run --rm -it -v $PWD:/go/src/github.com/ffimnsr/trader -p 8000:8000 golang:1.10.1 bash
> go get -u github.com/golang/dep/cmd/dep
> dep ensure
> go run main.go init_*.go
~~~
