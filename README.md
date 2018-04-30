# untilMongod

The `untilMongod` command attempts to connect to a MongoDB® server (`mongod`) or 
sharding front (`mongos`) for a maximum duration.

Its main use is as a way to start using a connection as soon as its server 
becomes available, without relying on manually adjusted timeouts.

[![GoDoc](https://godoc.org/github.com/fgm/untilMongod?status.svg)](https://godoc.org/github.com/fgm/untilMongod)
[![Build Status](https://travis-ci.org/fgm/untilMongod.svg?branch=master)](https://travis-ci.org/fgm/untilMongod)
[![Maintainability](https://api.codeclimate.com/v1/badges/84de4f16f20af011cee0/maintainability)](https://codeclimate.com/github/fgm/untilMongod/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/84de4f16f20af011cee0/test_coverage)](https://codeclimate.com/github/fgm/untilMongod/test_coverage)


## Syntax

    untilMongo -url mongodb://example.com:11117 -timeout 60
    
* `-url` is a typical MongoDB URL, defaulting to `mongodb://localhost:27017`
* `-timeout` is the maximum delay the command will wait before aborting.
* `-v` will increase verbosity, outputting messages on stderr on each retry.


## Exit codes

* 0: connection succeeded
* 1: connection could not be established, normal situation otherwise
* 2: another type of error occurred


## Installing the command

Assuming a Go 1.10 toolchain is available, just do:

```bash
go get github.com/fgm/untilMongod
```

This will install the `untilMongod` command in `$GOPATH/bin`, by default:

* `$HOME/go/bin` on UNIX-like systems (Linux, macOS) 
* `%USERPROFILE%\go\bin` on Windows systems.


## Running the Node.JS example

The `examples/example.bash` script show how to use `untilMongo` to run a Node.JS 
application only once the MongoDB server it tries to connect to has become 
available. The application itself is just an example showing the equivalent of
the `rs.Status()` mongo shell command.

To run it:

```bash
cd examples
npm i
bash example.bash
``` 


## Running tests

untilMongod uses the standard go testing package. The recommended way to run it
is:

```bash
go test -coverprofile coverage -covermode count ./...
go tool cover -html coverage
```


## IP information

* © 2018 Frederic G. MARAND
* Published under the [General Public License](LICENSE), version 3 or later.
* MongoDB is a trademark of MongoDB, Inc.
