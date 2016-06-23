package main

import (
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/op/go-logging"

    "github.com/BattleRattle/sexy/sentry"
    "github.com/BattleRattle/sexy/version"
    "github.com/BattleRattle/sexy/udp"
    "github.com/BurntSushi/toml"
)

type Config struct {
    UdpAddress  string
    SentryUrl   string
    Concurrency uint
    Buffer      uint
    LogLevel    string
    LogFile     string
}

var (
    configFile = flag.String("c", "/etc/sexy/sexy.toml", "Path to config file")
    showVersion = flag.Bool("version", false, "Show version information")

    format = logging.MustStringFormatter(`%{time:15:04:05.000} %{level:.4s} %{message}`)
    logger = logging.MustGetLogger("sexy")
)

func main() {
    flag.Parse()

    if *showVersion {
        printVersion()
        return
    }

    var cfg Config
    if _, err := toml.DecodeFile(*configFile, &cfg); err != nil {
        fmt.Fprintln(os.Stderr, "Failed to load config file: ", err)
        os.Exit(1)
    }

    logFile, err := os.OpenFile(cfg.LogFile, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0660)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to open log file " + cfg.LogFile, err)
        os.Exit(1)
    }

    lvl, err := logging.LogLevel(strings.ToUpper(cfg.LogLevel))
    if err != nil {
        fmt.Fprintln(os.Stderr, "Invalid log level. Allowed levels are: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL")
        os.Exit(1)
    }

    logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(logFile, "", 0), format))
    logging.SetLevel(lvl, "sexy")

    fmt.Println(fmt.Sprintf("SEXY - Sentry Proxy %s (%s)", version.Version, version.CommitHash))

    if cfg.UdpAddress == "" {
        fmt.Fprintln(os.Stderr, "No UDP address given to listen to\nRun `sexy -help` to display available arguments")
        os.Exit(1)
    }

    if cfg.SentryUrl == "" {
        fmt.Fprintln(os.Stderr, "No Sentry URL given\nRun `sexy -help` to display available arguments")
        os.Exit(1)
    }

    if cfg.Concurrency < 1 {
        fmt.Fprintln(os.Stderr, "Concurrency level must be 1 or higher")
        os.Exit(1)
    }

    if cfg.Buffer < 1 {
        fmt.Fprintln(os.Stderr, "Buffer size must be 1 or higher")
        os.Exit(1)
    }

    chMsg := make(chan sentry.Message, cfg.Buffer)
    defer close(chMsg)

    for i := uint(0); i < cfg.Concurrency; i++ {
        go sentry.NewWorker(cfg.SentryUrl, chMsg, logger).Run()
    }

    udp.NewServer(cfg.UdpAddress, chMsg, logger).Run()
}

func printVersion() {
    fmt.Println("SEXY - Sentry Proxy")
    fmt.Println()
    fmt.Println("Version:         ", version.Version)
    fmt.Println("Git Commit Hash: ", version.CommitHash)
    fmt.Println("Build Time:      ", version.BuildTime)
}