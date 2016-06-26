package udp

import (
    "strings"
    "errors"
    "net"

    "github.com/op/go-logging"

    "github.com/BattleRattle/sexy/sentry"
)

var (
    errUdpSig = errors.New(`Invalid UDP message (does not begin with "Sentry")`)
    errInvalidMessageFormat = errors.New("Invalid message format (expected AUTH_HEADER <NL> <NL> BASE64_BODY)")
)

type Server struct {
    Address string
    MsgChan chan <- sentry.Message
    Logger  *logging.Logger
}

func NewServer(addr string, msgChan chan <- sentry.Message, logger *logging.Logger) *Server {
    return &Server{Address: addr, MsgChan: msgChan, Logger: logger}
}

func (s *Server) Run() {
    resolvedAddr, err := net.ResolveUDPAddr("udp", s.Address)
    if err != nil {
        s.Logger.Fatalf("Resolving UDP address failed: %s", err)
    }

    conn, err := net.ListenUDP("udp", resolvedAddr)
    if err != nil {
        s.Logger.Fatalf("Opening UDP port failed: %s", err)
    }

    s.Logger.Infof("Listening on %s:%d (udp)", resolvedAddr.IP, resolvedAddr.Port)
    defer conn.Close()

    buf := make([]byte, 2 << 16)
    chUdp := make(chan string, 100)
    defer close(chUdp)

    go s.handleUdp(chUdp, s.MsgChan)

    for {
        n, remoteAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            s.Logger.Warningf("Reading UDP message failed: %s", err)
            continue
        }

        chUdp <- string(buf[0:n])
        s.Logger.Debugf("Received UDP package with %d bytes from %s", n, remoteAddr)
    }
}

// Handle UDP messages
func (s *Server) handleUdp(udpQueue <- chan string, sentryQueue chan <- sentry.Message) {
    var rawMsg string

    for {
        rawMsg = <-udpQueue

        msg, err := s.parseMessageUdp(&rawMsg)

        if err != nil {
            s.Logger.Info(err)
            continue
        }

        select {
        case sentryQueue <- *msg:
        default:
            s.Logger.Debug("Buffer is full, message will be dropped")
        }
    }
}

func (s *Server) parseMessageUdp(rawMsg *string) (*sentry.Message, error) {
    if !strings.HasPrefix(*rawMsg, "Sentry ") {
        return nil, errUdpSig
    }

    parts := strings.Split(*rawMsg, "\n\n")

    if len(parts) != 2 {
        return nil, errInvalidMessageFormat
    }

    return &sentry.Message{Header: parts[0], Body: parts[1]}, nil
}