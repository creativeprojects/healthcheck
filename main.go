package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	defaultTimeout = "10s"
)

func main() {
	err := healthcheck()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail: %s\n", err)
		os.Exit(1)
	}
}

func healthcheck() error {
	var rawTimeout, method string
	flag.StringVar(&rawTimeout, "timeout", defaultTimeout, "timeout for the healthcheck")
	flag.StringVar(&method, "method", http.MethodGet, "HTTP method for the healthcheck")
	flag.Parse()

	if len(flag.Args()) == 0 {
		return fmt.Errorf("missing url argument")
	}

	timeout, err := time.ParseDuration(rawTimeout)
	if err != nil {
		timeout, _ = time.ParseDuration(defaultTimeout)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, method, flag.Args()[0], http.NoBody)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("healthcheck failed with status code %d", resp.StatusCode)
	}

	return nil
}
