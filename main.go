package main

import (
	"context"
	"errors"
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

	cfg, err := ParseConfig(os.Environ())
	if err != nil {
		log.Fatalf("error parsing extension configuration: %v", err)
	}

	if cfg != nil {
		if err := gomplate.RunTemplates(cfg); err != nil {
			log.Fatal(err)
		}
		log.Printf("initialization successful... version=%s commit=%s date=%s", version, commit, date)
	} else {
		log.Println("initialization skipped: no templates found in configuration...")
	}

	if err := processEvents(ctx, client); err != nil {
		log.Fatal(err)
	}
}

type template struct {
	input  string
	output string
}

// ParseConfig generates a gomplate configuration using lambda function environment varialbes
func ParseConfig(envs []string) (*gomplate.Config, error) {
	cfg := &gomplate.Config{}

	// extract datasources & templates from envs
	var anonymousInput, anonymousOutput string
	templates := map[string]*template{}
	for _, env := range envs {
		switch {
		case strings.HasPrefix(env, "GOMPLATE_DATASOURCE_"):
			cfg.DataSources = append(cfg.DataSources, strings.TrimPrefix(env, "GOMPLATE_DATASOURCE_"))
		case strings.HasPrefix(env, "GOMPLATE_INPUT_"):
			pair := strings.SplitN(env, "=", 2)
			if len(pair) == 2 {
				name := strings.TrimPrefix(pair[0], "GOMPLATE_INPUT_")
				tpl, ok := templates[name]
				if !ok {
					tpl = &template{}
					templates[name] = tpl
				}
				tpl.input = pair[1]
			}
		case strings.HasPrefix(env, "GOMPLATE_INPUT="):
			anonymousInput = strings.TrimPrefix(env, "GOMPLATE_INPUT=")
		case strings.HasPrefix(env, "GOMPLATE_OUTPUT_"):
			pair := strings.SplitN(env, "=", 2)
			if len(pair) == 2 {
				name := strings.TrimPrefix(pair[0], "GOMPLATE_OUTPUT_")
				tpl, ok := templates[name]
				if !ok {
					tpl = &template{}
					templates[name] = tpl
				}
				tpl.output = pair[1]
			}
		case strings.HasPrefix(env, "GOMPLATE_OUTPUT="):
			anonymousOutput = strings.TrimPrefix(env, "GOMPLATE_OUTPUT=")
		}
	}

	if (anonymousInput != "" || anonymousOutput != "") && len(templates) > 0 {
		return nil, errors.New("only a single anonymous inline template or one or more named templates can be specified, not both")
	}
	if (anonymousInput != "") != (anonymousOutput != "") {
		return nil, errors.New("a single anonymous inline template requires both an input and output to be specified")
	}

	// validate templates
	if len(templates) > 0 {
		for name, tpl := range templates {
			if tpl == nil || len(tpl.input) == 0 || len(tpl.output) == 0 {
				return nil, fmt.Errorf("incomplete template detected: %s", name)
			}
			cfg.InputFiles = append(cfg.InputFiles, tpl.input)
			cfg.OutputFiles = append(cfg.OutputFiles, tpl.output)
		}
		return cfg, nil
	} else if (anonymousInput != "") && (anonymousOutput != "") {
		if stat, err := os.Stat(anonymousInput); err == nil && stat != nil && !stat.IsDir() {
			cfg.InputFiles = []string{anonymousInput}
		} else {
			cfg.Input = anonymousInput
		}
		cfg.OutputFiles = []string{anonymousOutput}
		return cfg, nil
	}

	return nil, nil
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
