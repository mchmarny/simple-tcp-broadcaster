# simple-tcp-broadcaster [![Build Status](https://travis-ci.org/mchmarny/simple-tcp-broadcaster.svg?branch=master)](https://travis-ci.org/mchmarny/simple-tcp-broadcaster)

The `simple-tcp-broadcaster` is a simple implementation of distributed client-server app which broadcasts messages to all connected clients. 

## Build

After you clone this repo, make sure to get the necessary dependencies, either by executing the `make deps` command or running acquiring the `dep` library and running it manually.

```shell
go get github.com/golang/dep/cmd/dep
dep ensure
```

After that you can build it using either the `make build` command or running the build manually.

```shell
go build -v -o ./bin/simple
```

## Usage

Running `simple-tcp-broadcaster` requires starting server and at least one client

### Start Server 

To start the server on port `5050` run:

```shell
./bin/simple server start --port 5050
```

### Start Client(s)

To start and connect the client to the already started server on the same machine run:

```shell
./bin/simple client connect --port 5050
```

If you connecting to remote machine you will also have to provide the address flag

```shell
./bin/simple client connect --address 10.0.0.10 --port 5050
```

Once running, the client listens for plain text commands and wraps them into a message structure and sends it to server where these messages are re-broadcasted to all connected clients.

### Message

THe messages look something like this, where the Data property caries your plain text message in binary format

```shell
    ID: 09b99dcc-0a06-4ffd-8c5f-d82921a890bf 
    ClientID: client-127-0-0-1-52982 
    CreatedAt: 2018-04-20 05:28:51.683622554 -0700 PDT 
    Status: 1 
    Data: [116 101 115 116 10]
```



