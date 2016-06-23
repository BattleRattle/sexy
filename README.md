Sexy (Sentry Proxy)
===================

**Sexy** is a UDP Proxy for [Sentry](https://getsentry.com/), written in Go.


Usage
-----

In order to run **Sexy** you just need to adjust your config file (`sexy.toml`) and run the binary.

```bash
sexy -c /path/to/sexy.toml
```

The following parameters need to be configured within the `sexy.toml` configuration:

```
# The local UDP address to listen on
udpAddress = "localhost:9001"

# The target Sentry base URL
sentryUrl = "https://sentry.example.org"

# The amount of workers to send out HTTP(S) requests in parallel
concurrency = 10

# The amount of Sentry messages to buffer (between receiving via UDP and sending out via HTTP(S))
buffer = 1000

# The logfile where some information (depending on log level) should be written to. The directory must already exist
logFile = "/var/log/sexy/sexy.log"

# The log level to be used to write into log file (debug, info, notice, warning, error, critical)
logLevel = "warning"
```

For more information about CLI arguments, run `sexy -help`

```bash
# sexy -help
Usage of sexy:
  -c string
    	Path to config file (default "/etc/sexy/sexy.toml")
  -version
    	Show version information
```


Contributing
------------

**Sexy** should only serve its purpose of an UDP proxy for Sentry. But in case you find any bug or think it can be improved in some way,
please don't hesitate to create an issue or - even better - create a pull request with your proposed solution.


License
-------

This project is released under the terms of the [MIT license] (https://github.com/BattleRattle/sexy/blob/master/LICENSE)