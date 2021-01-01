# untilMongod

The `untilMongod` command attempts to connect to a MongoDB® server (`mongod`) or 
sharding front (`mongos`) for a maximum duration.

Its main use is as a way to start using a connection as soon as its server 
becomes available, without relying on manually adjusted timeouts.

[![GoDoc](https://godoc.org/github.com/fgm/untilMongod?status.svg)](https://godoc.org/github.com/fgm/untilMongod)
[![Go Report Card](https://goreportcard.com/badge/github.com/fgm/untilMongod)](https://goreportcard.com/report/github.com/fgm/untilMongod)
[![Build Status](https://travis-ci.org/fgm/untilMongod.svg?branch=develop)](https://travis-ci.org/fgm/untilMongod)
[![Maintainability](https://api.codeclimate.com/v1/badges/84de4f16f20af011cee0/maintainability)](https://codeclimate.com/github/fgm/untilMongod/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/84de4f16f20af011cee0/test_coverage)](https://codeclimate.com/github/fgm/untilMongod/test_coverage)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ffgm%2FuntilMongod.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Ffgm%2FuntilMongod?ref=badge_shield)

## Syntax

    untilMongo -url mongodb://example.com:11117 -timeout 60
    
* `-url` is a typical MongoDB URL, defaulting to `mongodb://localhost:27017`
  * It supports the MGO driver [Dial() URL extensions], which is necessary to support connecting to servers directly without
    reaching for the replica set, and other situations like authentication or pool tuning. 
   * Specifically, to connect to a server started as part of a yet-unconfigured replica set, the URL must contain a
     `connect=direct` query, like:
        
         untilMongo -url mongodb://example.com:11117?connect=direct -timeout 60
* `-timeout` is the maximum delay the command will wait before aborting.
* `-v` will increase verbosity, outputting messages on stderr on each retry.

[Dial() URL extensions]: https://godoc.org/github.com/globalsign/mgo#Dial


## Exit codes

* 0: connection succeeded
* 1: connection could not be established, normal situation otherwise
* 2: another type of error occurred


## Installing the command

Assuming a Go 1.15 toolchain is available, just do:

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

untilMongod uses the standard go testing package. Tests should pass whether or
not a `mongod|mongos` is available at the default `mongodb://localhost:27017`
URL, but they will only cover the successful connection path if an instance is
available.

The recommended ways to run tests are:

```bash
# Running all tests and generating coverage.
go test -coverprofile coverage -covermode count ./...
go tool cover -html coverage

# Running only unit tests
go test -short ./...

# Running only integration tests
go test -run 'Integration$' ./...
```


## IP information

* © 2018-2020 Frederic G. MARAND.
* Published under the [General Public License](LICENSE), version 3 or later.
* MongoDB is a trademark of MongoDB, Inc.


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ffgm%2FuntilMongod.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Ffgm%2FuntilMongod?ref=badge_large)