# untilMongod

The `untilMongod` command attempts to connect to a MongoDB server (`mongod`) or 
sharding front (`mongos`) for a maximum duration.

Its main use is as a way to start using a Mongo connection as soon as it becomes
available, without relying on manually adjusted timeouts.

## Syntax

    untilMongo -url mongodb://example.com:11117 -timeout 60
    
* `-url` is a typical MongoDB URL, defaulting to `mongodb://localhost:27017`
* `-timeout` is the maximum delay the command will wait before aborting.

## Exit codes

* 0: connection succeeded
* 1: connection could not be established, normal situation otherwise
* 2: another type of error occurred

## IP information

* &copy; 2018 Frederic G. MARAND
* Published under the [General Public License](LICENSE), version 3 or later.