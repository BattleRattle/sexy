package sentry

import (
    "errors"
    "fmt"
    "net/http"
    "net/url"
    "strings"

    "github.com/op/go-logging"

    "github.com/BattleRattle/sexy/version"
)

var (
    errInvalidScheme = errors.New(`Invalid Sentry URL: must use "http" or "https" as protocol`)
)

type Message struct {
    Header string
    Body   string
}

type Client struct {
    Protocol   string
    Host       string
    Logger     *logging.Logger
    httpClient *http.Client
}

func NewClient(baseUrl string, logger *logging.Logger) (client *Client, err error) {
    parsedUrl, err := url.Parse(baseUrl)
    if err != nil {
        return
    }

    if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
        return nil, errInvalidScheme
    }

    client = &Client{Protocol: parsedUrl.Scheme, Host: parsedUrl.Host, Logger: logger, httpClient: &http.Client{}}

    return
}

func (c *Client) Send(msg Message) (err error) {
    url := fmt.Sprintf("%s://%s/api/store/", c.Protocol, c.Host)

    req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(msg.Body))
    if err != nil {
        return err
    }

    req.Header.Add("X-Sentry-Auth", msg.Header)
    req.Header.Add("Content-Type", "application/octet-stream")
    req.Header.Add("User-Agent", "Sexy/" + version.Version)

    c.Logger.Debug("Sending message to Sentry")
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        c.Logger.Debug("Successfully delivered message to Sentry")
    } else {
        c.Logger.Warningf("Invalid response (status code %d)", resp.StatusCode)
    }

    return
}