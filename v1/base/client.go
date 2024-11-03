package base

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Client struct {
	httpClient *http.Client
	logger     *log.Logger
	verbose    bool
}

func NewClient(verbose bool) (*Client, error) {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	httpClient := &Client{
		httpClient: client,
		logger:     logger,
		verbose:    verbose,
	}

	return httpClient, nil
}

func (c *Client) log(format string, v ...interface{}) {
	if c.verbose {
		c.logger.Printf(format, v...)
	}
}

func (c *Client) Do(ctx context.Context, r *Request) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, r.Method, r.Url, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header
	c.log("Request : %s %s %v %v", req.Method, req.URL, req.Header, req.Body)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		// TODO add more handling
		return nil, errors.New(res.Status)
	}

	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the returned error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
