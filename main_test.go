package main_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	extension "github.com/cludden/gomplate-lambda-extension"
	"github.com/hairyhenderson/gomplate/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	dir := t.TempDir()

	cases := map[string]struct {
		environment func() ([]string, error)
		assert      func(*testing.T, bool, error)
	}{
		"err_both_named_and_anonymous": {
			environment: func() ([]string, error) {
				return []string{
					"GOMPLATE_INPUT=/tmp/foo.tpl.json",
					"GOMPLATE_OUTPUT=/tmp/foo.json",
					"GOMPLATE_INPUT_foo=/tmp/bar/foo.tpl.json",
					"GOMPLATE_OUTPUT_foo=/tmp/bar/foo.json",
				}, nil
			},
			assert: func(t *testing.T, ok bool, err error) {
				assert.False(t, ok)
				assert.Error(t, err)
			},
		},
		"err_anonymous_missing_input": {
			environment: func() ([]string, error) {
				return []string{
					"GOMPLATE_OUTPUT=/tmp/foo.json",
				}, nil
			},
			assert: func(t *testing.T, ok bool, err error) {
				assert.False(t, ok)
				assert.Error(t, err)
			},
		},
		"err_anonymous_missing_output": {
			environment: func() ([]string, error) {
				return []string{
					"GOMPLATE_INPUT=/tmp/foo.tpl.json",
				}, nil
			},
			assert: func(t *testing.T, ok bool, err error) {
				assert.False(t, ok)
				assert.Error(t, err)
			},
		},
		"err_incomplete_named": {
			environment: func() ([]string, error) {
				return []string{
					"GOMPLATE_INPUT_foo=/tmp/bar/foo.tpl.json",
					"GOMPLATE_OUTPUT_foos=/tmp/bar/foo.json",
				}, nil
			},
			assert: func(t *testing.T, ok bool, err error) {
				assert.False(t, ok)
				assert.Error(t, err)
			},
		},
		"envs_inline": {
			environment: func() ([]string, error) {
				return []string{
					`GOMPLATE_INPUT=foo: {{ env.Getenv "foo" "bar" }}`,
					fmt.Sprintf("GOMPLATE_OUTPUT=%s/foo", dir),
				}, nil
			},
			assert: func(t *testing.T, ok bool, err error) {
				assert.NoError(t, err)
				out, _ := ioutil.ReadFile(fmt.Sprintf("%s/foo", dir))
				assert.Equal(t, "foo: bar", string(out))
			},
		},
		"envs_file": {
			environment: func() ([]string, error) {
				if err := ioutil.WriteFile(fmt.Sprintf("%s/foo.tpl.yml", dir), []byte(`foo: {{ env.Getenv "foo" "bar" }}`), 0777); err != nil {
					return nil, err
				}
				return []string{
					fmt.Sprintf("GOMPLATE_INPUT=%s/foo.tpl.yml", dir),
					fmt.Sprintf("GOMPLATE_OUTPUT=%s/foo.yml", dir),
				}, nil
			},
			assert: func(t *testing.T, ok bool, err error) {
				assert.NoError(t, err)
				out, _ := ioutil.ReadFile(fmt.Sprintf("%s/foo.yml", dir))
				assert.Equal(t, "foo: bar", string(out))
			},
		},
		"envs_named": {
			environment: func() ([]string, error) {
				if err := ioutil.WriteFile(fmt.Sprintf("%s/bar.tpl.yml", dir), []byte(`foo: {{ env.Getenv "foo" "bar" }}`), 0777); err != nil {
					return nil, err
				}
				return []string{
					fmt.Sprintf("GOMPLATE_INPUT_bar=%s/bar.tpl.yml", dir),
					fmt.Sprintf("GOMPLATE_OUTPUT_bar=%s/bar.yml", dir),
				}, nil
			},
			assert: func(t *testing.T, ok bool, err error) {
				assert.NoError(t, err)
				out, _ := ioutil.ReadFile(fmt.Sprintf("%s/bar.yml", dir))
				assert.Equal(t, "foo: bar", string(out))
			},
		},
		"files_inline": {
			environment: func() ([]string, error) {
				f := fmt.Sprintf("%s/files.config.json", dir)
				if err := ioutil.WriteFile(f, []byte(`{"foo":"bar"}`), 0777); err != nil {
					return nil, fmt.Errorf("error writing config file: %v", err)
				}
				return []string{
					`GOMPLATE_INPUT=foo: {{ (ds "foo").foo }}`,
					fmt.Sprintf("GOMPLATE_OUTPUT=%s/foo", dir),
					fmt.Sprintf("GOMPLATE_DATASOURCE_foo=file://%s", f),
				}, nil
			},
			assert: func(t *testing.T, ok bool, err error) {
				assert.NoError(t, err)
				out, _ := ioutil.ReadFile(fmt.Sprintf("%s/foo", dir))
				assert.Equal(t, "foo: bar", string(out))
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			envs, err := c.environment()
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			cfg, err := extension.ParseConfig(envs)
			if err != nil {
				c.assert(t, false, err)
				return
			}
			log.Printf("%+v\n", cfg)

			err = gomplate.RunTemplates(cfg)
			c.assert(t, true, err)
		})
	}
}
