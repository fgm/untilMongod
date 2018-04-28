# untilMongod

The `untilMongod` command attempts to connect to a MongoDB &reg; server (`mongod`) or 
sharding front (`mongos`) for a maximum duration.

Its main use is as a way to start using a connection as soon as its server 
becomes available, without relying on manually adjusted timeouts.


## Syntax

    untilMongo -url mongodb://example.com:11117 -timeout 60
    
* `-url` is a typical MongoDB URL, defaulting to `mongodb://localhost:27017`
* `-timeout` is the maximum delay the command will wait before aborting.


## Exit codes

* 0: connection succeeded
* 1: connection could not be established, normal situation otherwise
* 2: another type of error occurred


## Usage example

The `examples/example.bash` script show how to use `untilMongo` to run a NodeJS 
application only once the MongoDB server it tries to connect to has become 
available.

To run it:

```bash
cd examples
npm i
bash example.bash
``` 


## IP information

* &copy; 2018 Frederic G. MARAND
* Published under the [General Public License](LICENSE), version 3 or later.
* MongoDB is a trademark of MongoDB, Inc.