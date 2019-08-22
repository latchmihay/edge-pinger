package config

import (
	"github.com/hashicorp/go-multierror"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"net"
	"time"

	"github.com/hashicorp/hcl"
	"github.com/pkg/errors"
)

type Config struct {
	Count     int      `hcl:"count"`
	Timeout   string   `hcl:"timeout"`
	Interval  string   `hcl:"interval"`
	Addresses []string `hcl:"addresses"`
}

func LoadConfigFile(filepath string) (string, error) {
	hclText, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(hclText), nil
}

var multiErrors *multierror.Error

func ParseConfig(hclText string) (*Config, error) {
	result := &Config{}

	hclParseTree, err := hcl.Parse(hclText)
	if err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&result, hclParseTree); err != nil {
		return nil, err
	}

	if result.Timeout != "" {
		_, err = time.ParseDuration(result.Timeout)
		if err != nil {
			err = errors.Wrap(err, "Timeout (format: 15s)")
			multiErrors = multierror.Append(multiErrors, err)
		}
	}

	if result.Interval != "" {
		_, err = time.ParseDuration(result.Interval)
		if err != nil {
			err = errors.Wrap(err, "Interval (format: 15s)")
			multiErrors = multierror.Append(multiErrors, err)
		}
	}

	result.Addresses = funk.UniqString(result.Addresses)
	for _, addr := range result.Addresses {
		_, err := net.ResolveIPAddr("ip4:icmp", addr)
		if err != nil {
			multiErrors = multierror.Append(multiErrors, err)
		}
	}

	return result, multiErrors.ErrorOrNil()
}
