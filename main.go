package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "strings"
    "syscall"

    "github.com/BurntSushi/toml"
    "github.com/op/go-logging"

    "github.com/BattleRattle/sexy/log"
    "github.com/BattleRattle/sexy/sentry"
    "github.com/BattleRattle/sexy/udp"
    "github.com/BattleRattle/sexy/version"
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

    format = logging.MustStringFormatter(`%{time:2006-01-02 15:04:05.000} %{level:.4s} %{message}`)
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

    logWriter, err := log.NewFileWriter(cfg.LogFile, 0644)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to open log file " + cfg.LogFile, err)
        os.Exit(1)
    }

    lvl, err := logging.LogLevel(strings.ToUpper(cfg.LogLevel))
    if err != nil {
        fmt.Fprintln(os.Stderr, "Invalid log level. Allowed levels are: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL")
        os.Exit(1)
    }

    logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(logWriter, "", 0), format))
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

    go runSignalHandler(logWriter)

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

func runSignalHandler(logWriter *log.FileWriter) {
    chSig := make(chan os.Signal)
    defer close(chSig)

    signal.Notify(chSig, syscall.SIGUSR1)
    logger.Debug("Registered Signal Handler")

    for sig := range chSig {
        logger.Infof("Received %s", sig)

        if err := logWriter.Reopen(); err != nil {
            logger.Warningf("Unable to reopen log file: %s", err)
        }
    }
}