package main_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	extension "github.com/cludden/gomplate-lambda-extension"
	"github.com/hairyhenderson/gomplate/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	dir := t.TempDir()

	cases := map[string]struct {
		input       string
		datasources func() ([]string, error)
		assert      func(*testing.T, string)
	}{
		"envs": {
			input: `foo: {{ env.Getenv "foo" "bar" }}`,
			assert: func(t *testing.T, out string) {
				assert.Equal(t, "foo: bar", out)
			},
		},
		"files": {
			input: `foo: {{ (ds "foo").foo }}`,
			datasources: func() ([]string, error) {
				f := fmt.Sprintf("%s/files.config.json", dir)
				if err := ioutil.WriteFile(f, []byte(`{"foo":"bar"}`), 0777); err != nil {
					return nil, fmt.Errorf("error writing config file: %v", err)
				}
				return []string{fmt.Sprintf("GOMPLATE_DATASOURCE_foo=file://%s", f)}, nil
			},
			assert: func(t *testing.T, out string) {
				assert.Equal(t, "foo: bar", out)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			out := fmt.Sprintf("%s/%s", dir, name)
			var ds []string
			var err error
			if c.datasources != nil {
				ds, err = c.datasources()
				if err != nil {
					t.Fatalf("error initializing datasources: %v", err)
				}
			}
			cfg, err := extension.ParseConfig(c.input, out, ds)
			if !assert.NoError(t, err, "unexpected error parsing config") {
				t.FailNow()
			}
			err = gomplate.RunTemplates(cfg)
			if !assert.NoError(t, err, "unexpected error executing gomplate") {
				t.Fail()
			}
			rendered, err := ioutil.ReadFile(out)
			if !assert.NoError(t, err, "unexpected error reading output") {
				t.Fail()
			}
			c.assert(t, string(rendered))
		})
	}
}
