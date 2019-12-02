// periodically crawl metrics from http endpoint and send to kafka
package kree

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

// Message is used to communicate between Collector and Publisher
type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

// Collector is for periodically reading from metric endpoints
// and send message to Publisher through its channel
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
	log.Println("Starting collector...")
	ctx := context.Background()
	ticker := time.NewTicker(60 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				endpoints, err := c.repository.GetAll()
				if err != nil {
					log.Fatalf("Could not get endpoint %v", err)
				}
				if endpoints == nil {
					log.Println("No endpoints found")
				}
				for _, e := range endpoints {
					c.input <- e
				}
			case <-c.quit:
				ticker.Stop()
				log.Println("Received shutting down request")
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
			err := c.makeRequest(ctx, e)
			if err != nil {
				log.Fatalf("Could not collect metrics from %v, err: %v", e, err)
			}
		case <-c.quit:
			break
		}
	}
}

func (c *Collector) makeRequest(ctx context.Context, e *endpoint) error {
	log.Printf("Making request to endpoint %v", e)
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
