Sexy (Sentry Proxy)
===================

**Sexy** is a UDP Proxy for [Sentry](https://getsentry.com/), written in Go.


Usage
-----

In order to run **Sexy** you just need to provide a local UDP address (`-u`) and the base URL of the sentry instance (`-s`).

```bash
sexy -u localhost:9001 -s https://sentry.example.org
```

For more information about CLI arguments, run `sexy -help`

```bash
# sexy -help
Usage of sexy:
  -b uint
    	Buffer size for pending requests (default 1000)
  -c uint
    	Concurrency level for passing requests to Sentry (default 1)
  -loglevel string
    	The log level, any of: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL (default "WARNING")
  -s string
    	Sentry base URL (e.g. https://sentry.example.org)
  -u string
    	UDP address to listen on (default "localhost:9001")
  -version
    	Show version information
```

