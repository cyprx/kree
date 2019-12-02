// periodically crawl metrics from http endpoint and send to kafka
package kree

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

type Collector struct {
	topic      string
	repository EndpointRepository
	parser     *parser
	c          *http.Client
	input      chan *endpoint
	output     chan *Message
	quit       chan os.Signal
}

func NewCollector(r EndpointRepository, o chan *Message, q chan os.Signal) *Collector {
	ch := make(chan *endpoint, 1)
	return &Collector{
		repository: r,
		parser:     newParser(o),
		c:          &http.Client{},
		input:      ch,
		output:     o,
		quit:       q,
	}
}

func (c *Collector) RegisterEndpoint(e *endpoint) {
}

func (c *Collector) Run() {
	ctx := context.Background()
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				endpoints, err := c.repository.GetAll()
				if err != nil {
					log.Fatalf("could not get endpoint %v", err)
				}
				for _, e := range endpoints {
					c.input <- e
				}
			case <-c.quit:
				ticker.Stop()
				close(c.input)
			}
		}
	}()

	c.collect(ctx)
}

func (c *Collector) collect(ctx context.Context) {
	for {
		select {
		case e := <-c.input:
			c.makeRequest(ctx, e)
		case <-c.quit:
			break
		}
	}
}

func (c *Collector) makeRequest(ctx context.Context, e *endpoint) error {
	req, err := http.NewRequest("GET", e.Url(), nil)

	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return c.parser.parseTextToMessage(resp.Body)
}
