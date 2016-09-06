// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"
	"regexp"
)

type Config struct {
	Period time.Duration `config:"period"`
        Path string `config:"path"`
	IncludeLines []*regexp.Regexp `config:"include_lines"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}
