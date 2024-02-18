# untilMongod

The `untilMongod` command attempts to connect to a MongoDB® server (`mongod`) or 
sharding front (`mongos`) for a maximum duration.

Its main use is as a way to start using a connection as soon as its server 
becomes available, without relying on manually adjusted timeouts.

[![GoDoc](https://godoc.org/github.com/fgm/untilMongod?status.svg)](https://godoc.org/github.com/fgm/untilMongod)
[![Go Report Card](https://goreportcard.com/badge/github.com/fgm/untilMongod)](https://goreportcard.com/report/github.com/fgm/untilMongod)
[![Build Status](https://app.travis-ci.com/fgm/untilMongod.svg?branch=main)](https://app.travis-ci.com/fgm/untilMongod)
[![Maintainability](https://api.codeclimate.com/v1/badges/84de4f16f20af011cee0/maintainability)](https://codeclimate.com/github/fgm/untilMongod/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/84de4f16f20af011cee0/test_coverage)](https://codeclimate.com/github/fgm/untilMongod/test_coverage)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/fgm/untilMongod/badge)](https://securityscorecards.dev/viewer/?uri=github.com/fgm/untilMongod)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ffgm%2FuntilMongod.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2Ffgm%2FuntilMongod?ref=badge_small)

## Syntax

    untilMongo -url mongodb://example.com:11117 -timeout 60
    
* `-url` is a typical MongoDB URL, defaulting to `mongodb://localhost:27017`
  * It supports the MongoDB driver [Dial() URL options], which is necessary to support connecting to servers 
    directly without
    reaching for the replica set, and other situations like authentication or pool tuning. 
   * Specifically, to connect to a server started as part of a yet-unconfigured replica set, the URL must contain a
     `directConnection=true` query, like:
        
         untilMongo -url mongodb://example.com:11117?directConnection=true -timeout 60
* `-timeout` is the maximum delay (in seconds) the command will wait before aborting.
* `-v` will increase verbosity, outputting messages on stderr on each retry.

[Dial() URL options]: https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#std-label-golang
-connection-options


## Exit codes

* 0: connection succeeded
* 1: connection could not be established, normal situation otherwise
* 2: another type of error occurred


## Installing the command

Assuming a Go 1.22 or later toolchain is available, just do:

```bash
go install github.com/fgm/untilMongod@latest
```

This will install the `untilMongod` command in `$GOPATH/bin`, by default:

* `$HOME/go/bin` on UNIX-like systems (Linux, macOS) 
* `%USERPROFILE%\go\bin` on Windows systems.


## Running the Node.JS example

The `examples/example.bash` script show how to use `untilMongo` to run a Node.JS 
application only once the MongoDB server it tries to connect to has become 
available. The application itself is just an example showing the equivalent of
the `rs.Status()` mongo shell command.

It assumes that:

- MongoDB Community has been installed
- Go is installed on the host

To run it:

```bash
cd examples
npm i
bash example.bash
``` 


## Running tests

untilMongod uses the standard go testing package. 
Tests should pass, regardless of whether a `mongod|mongos` is available at the default `mongodb://localhost:27017` URL, 
but they will only cover the successful connection path if an instance is available.

The recommended ways to run tests are:

```bash
# Running all tests and generating coverage.
go test -race -coverprofile coverage -covermode atomic ./...
go tool cover -html coverage

# Running only unit tests
go test -race -short ./...

# Running only integration tests
go test -race -run 'Integration$' ./...
```


## IP information

* © 2018-2024 Frederic G. MARAND.
* Published under the [General Public License](LICENSE), version 3 or later.
* MongoDB is a trademark of MongoDB, Inc.


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ffgm%2FuntilMongod.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Ffgm%2FuntilMongod?ref=badge_large)