package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/op/go-logging"

    "github.com/BattleRattle/sexy/sentry"
    "github.com/BattleRattle/sexy/version"
    "github.com/BattleRattle/sexy/udp"
    "strings"
)

var (
    udpAddress = flag.String("u", "localhost:9001", "UDP address to listen on")
    sentryUrl = flag.String("s", "", "Sentry base URL (e.g. https://sentry.example.org)")
    concurrency = flag.Uint("c", 1, "Concurrency level for passing requests to Sentry")
    buffer = flag.Uint("b", 1000, "Buffer size for pending requests")
    logLevel = flag.String("loglevel", "WARNING", "The log level, any of: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL")
    showVersion = flag.Bool("version", false, "Show version information")

    format = logging.MustStringFormatter(`%{color}%{time:15:04:05} %{level:.4s}%{color:reset} %{message}`)
    logBackend = logging.NewLogBackend(os.Stdout, "", 0)
    logger = logging.MustGetLogger("sexy")
)

func main() {
    flag.Parse()

    lvl, err := logging.LogLevel(strings.ToUpper(*logLevel))
    if err != nil {
        fmt.Fprintln(os.Stderr, "Invalid log level. Allowed levels are: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL")
        os.Exit(1)
    }

    logging.SetBackend(logging.NewBackendFormatter(logBackend, format))
    logging.SetLevel(lvl, "sexy")

    if *showVersion {
        printVersion()
        return
    }

    if *udpAddress == "" {
        fmt.Fprintln(os.Stderr, "No UDP address given to listen to\nRun `sexy -help` to display available arguments")
        os.Exit(1)
    }

    if *sentryUrl == "" {
        fmt.Fprintln(os.Stderr, "No Sentry URL given\nRun `sexy -help` to display available arguments")
        os.Exit(1)
    }

    if *concurrency < 1 {
        fmt.Fprintln(os.Stderr, "Concurrency level must be 1 or higher")
        os.Exit(1)
    }

    if *buffer < 1 {
        fmt.Fprintln(os.Stderr, "Buffer size must be 1 or higher")
        os.Exit(1)
    }

    chMsg := make(chan sentry.Message, *buffer)

    for i := uint(0); i < *concurrency; i++ {
        go sentry.NewWorker(*sentryUrl, chMsg, logger).Run()
    }

    udp.NewServer(*udpAddress, chMsg, logger).Run()
}

func printVersion() {
    fmt.Println("Sentry Proxy")
    fmt.Println()
    fmt.Println("Version:         ", version.Version)
    fmt.Println("Git Commit Hash: ", version.CommitHash)
    fmt.Println("Build Time:      ", version.BuildTime)
}