package sentry

import (
    "github.com/op/go-logging"
)

type Worker struct {
    SentryUrl string
    Queue     <-chan Message
    Logger    *logging.Logger
}

func NewWorker(sentryUrl string, queue <-chan Message, logger *logging.Logger) *Worker {
    return &Worker{SentryUrl: sentryUrl, Queue: queue, Logger: logger}
}

func (w *Worker) Run() {
    client, err := NewClient(w.SentryUrl, w.Logger)
    if err != nil {
        w.Logger.Fatal(err)
    }

    for {
        if err := client.Send(<-w.Queue); err != nil {
            w.Logger.Warningf("An error occured while sending message to Sentry: %s", err)
        }
    }
}