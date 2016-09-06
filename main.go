package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/tak7iji/pfrbeat/beater"
)

func main() {
	err := beat.Run("pfrbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
