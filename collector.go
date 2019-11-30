// periodically crawl metrics from http endpoint and send to kafka
package kree

import (
	"context"
	"log"
	"net/http"
)

type Collector struct {
	repository EndpointRepository
	parser     *parser
	c          *http.Client
}

func NewCollector() *Collector {
	return &Collector{
		repository: NewEndpointRepository(),
		parser:     newParser(),
		c:          &http.Client{},
	}
}

func (c *Collector) RegisterEndpoint(e *endpoint) {
}

func (c *Collector) Collect() {
	endpoints, err := c.repository.GetAll()

	if err != nil {
	}

	for _, endpoint := range endpoints {
		log.Println(endpoint)
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

	json, err := c.parser.parseTextToJSON(resp.Body)
	log.Println(string(json))

	return err

}
