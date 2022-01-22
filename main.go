package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cludden/gomplate-lambda-extension/internal/extension"
	"github.com/hairyhenderson/gomplate/v3"
)

const (
	extensionName = "gomplate-lambda-extension"
)

// build variables
var (
	version string = "development"
	commit  string = "development"
	date    string = time.Now().UTC().Format(time.RFC3339)
)

// runtime configuration variables
var (
	Input  = os.Getenv("GOMPLATE_INPUT")
	Output = os.Getenv("GOMPLATE_OUTPUT")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	client := extension.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigs
		cancel()
	}()

	if _, err := client.Register(ctx, extensionName); err != nil {
		log.Fatalf("error registering extension: %v", err)
	}

	cfg, err := ParseConfig(Input, Output, os.Environ())
	if err != nil {
		log.Fatalf("error parsing extension configuration: %v", err)
	}

	if err := gomplate.RunTemplates(cfg); err != nil {
		log.Fatal(err)
	}
	log.Printf("initialization successful... version=%s commit=%s date=%s", version, commit, date)

	if err := processEvents(ctx, client); err != nil {
		log.Fatal(err)
	}
}

// ParseConfig generates a gomplate configuration using lambda function environment varialbes
func ParseConfig(input, output string, envs []string) (*gomplate.Config, error) {
	cfg := &gomplate.Config{
		Input:       input,
		OutputFiles: []string{output},
		DataSources: []string{},
	}
	for _, env := range envs {
		if strings.HasPrefix(env, "GOMPLATE_DATASOURCE_") {
			cfg.DataSources = append(cfg.DataSources, strings.TrimPrefix(env, "GOMPLATE_DATASOURCE_"))
		}
	}

	missing := []string{}
	if len(cfg.Input) == 0 {
		missing = append(missing, "GOMPLATE_INPUT")
	}
	if len(cfg.OutputFiles[0]) == 0 {
		missing = append(missing, "GOMPLATE_OUTPUT")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variable(s): %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

// processEvents blocks until shutdown event is received or cancelled via os signals
func processEvents(ctx context.Context, client *extension.Client) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			res, err := client.NextEvent(ctx)
			if err != nil {
				return fmt.Errorf("error polling next event: %v", err)
			}
			if res.EventType == extension.Shutdown {
				return nil
			}
		}
	}
}
