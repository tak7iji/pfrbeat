package beater

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/tak7iji/pfrbeat/config"
)

type Pfrbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Pfrbeat{
		done: make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Pfrbeat) Run(b *beat.Beat) error {
	logp.Info("pfrbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

                bt.readLine(bt.config.Path, b.Name)

		logp.Info("Event sent")
	}
}

func (bt *Pfrbeat) readLine(filename string, beatname string) {
	var fp *os.File

	fp, _ = os.Open(filename)
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		if line := scanner.Text(); bt.isIncludeLines(line) {
			bt.publish(line, beatname)
		}
	}

	// omit error handling

}

func (bt *Pfrbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Pfrbeat) isIncludeLines(line string) bool {
    if len(bt.config.IncludeLines) > 0 {
        for _, rexp := range bt.config.IncludeLines {
            if rexp.MatchString(line) {
                return true
            }
        }
    }
    return false
}

func (bt *Pfrbeat) publish(line string, beatname string) {
	event := common.MapStr {
		"@timestamp": common.Time(time.Now()),
		"type":	beatname,
		"message": line,
	}

	bt.client.PublishEvent(event)
}
