package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	expected = &Config{
		Count:     5,
		Timeout:   "5s",
		Interval:  "60s",
		Addresses: []string{"google.com", "facebook.com", "github.com"},
	}
	testConfig = `count = 5
		timeout = "5s"
		interval = "60s"
		addresses = [
			"google.com",
			"facebook.com",
			"facebook.com",
			"github.com",
		]`
	badAddresses = `addresses = [
			"google.co1m",
			"facebook.com1",
			"github.com1",
		]`
)

func TestConfigParsing(t *testing.T) {
	require := require.New(t)

	config, err := ParseConfig(testConfig)
	require.NoError(err)
	require.Equal(expected, config)
}

func TestBadAddresses(t *testing.T) {
	require := require.New(t)

	_, err := ParseConfig(badAddresses)
	t.Log(err)
	require.Error(err)
}
