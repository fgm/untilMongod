# untilMongod

> [!IMPORTANT]
> **This repository is archived and no longer maintained.**
> Use [wait4x](https://github.com/wait4x/wait4x) instead:
> its `mongodb` command performs the same check,
> and it is maintained, packaged for most platforms,
> and covers many other services.

The `untilMongod` command waited until a MongoDB® server (`mongod`)
or sharding front (`mongos`) accepted connections, for a maximum duration.


## Migrating to wait4x

Install wait4x from any of the channels listed in [its README](https://github.com/wait4x/wait4x#readme),
e.g. `brew install wait4x`, `apk add wait4x`, the `wait4x/wait4x` Docker image,
or `go install wait4x.dev/v3/cmd/wait4x@latest`.

Then replace:

```bash
untilMongod -url 'mongodb://example.com:11117/?directConnection=true' -timeout 60
```

with:

```bash
wait4x mongodb 'mongodb://example.com:11117/?directConnection=true' -t 60s
```

### Differences to check in your scripts

| Topic                 | untilMongod                         | wait4x                                         |
|-----------------------|-------------------------------------|------------------------------------------------|
| URL                   | `-url` flag, default `mongodb://localhost:27017` | positional argument, required      |
| Timeout               | `-timeout 30`: seconds, default 30  | `-t 30s`: a duration, **default 10s**          |
| Output                | silent unless `-v`                  | logs each attempt and the final error unless `-q` |
| `-v`                  | verbose                             | **invert check: succeeds when the server is down** |
| Exit code on success  | 0                                   | 0                                              |
| Exit code on timeout  | 1                                   | **124**                                        |
| Exit code on other errors | 2                               | **1**                                          |
| Authentication failure | exits 2 immediately                | retries until the timeout, then exits 124      |

- **Do not carry `-v` over.**
  In wait4x it inverts the check,
  so a script that kept it would proceed precisely when MongoDB is unavailable.
- **Exit code 1 changes meaning**, from "timed out" to "other error".
  Scripts testing for a specific code need updating.
- The URL options described in the [MongoDB Go driver connection options] still apply,
  including `directConnection=true` to reach a member of a replica set that is not yet initiated.

wait4x also supports:

- checking several URLs in parallel,
- running a command once the check succeeds: `wait4x mongodb 'mongodb://...' -- ./start-app.sh`,
- exponential backoff tuning, with `--backoff-policy exponential`.

[MongoDB Go driver connection options]: https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#connection-options


## Existing installations

The last release remains installable, but will receive no fixes:

```bash
go install github.com/fgm/untilMongod@latest
```

Its syntax was:

    untilMongod [-url mongodb://localhost:27017] [-timeout 30] [-v]


## IP and Licensing information

* © 2018-2026 Frederic G. MARAND.
* Published under the [General Public License](LICENSE), version 3.0 or later (SPDX: GPL-3.0-or-later)
* MongoDB is a trademark of MongoDB, Inc.
