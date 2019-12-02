package kree

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/prom2json"
)

type Metric struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	Labels string      `json:"labels"`
	Value  interface{} `json:"value"`
}

type parser struct {
	p  expfmt.TextParser
	ch chan *Message
}

func newParser(ch chan *Message) *parser {
	return &parser{expfmt.TextParser{}, ch}
}

func (p *parser) parseTextToMessage(line io.Reader) error {
	mfs, err := p.p.TextToMetricFamilies(line)
	if err != nil {
		return err
	}
	for _, mf := range mfs {
		f := prom2json.NewFamily(mf)

		// TODO: handle type histogram
		if f.Type == dto.MetricType_HISTOGRAM.String() {
			continue
		}
		return p.parseMetric(f)
	}
	return nil
}

func (p *parser) parseMetric(f *prom2json.Family) error {
	for _, fm := range f.Metrics {
		var (
			metric = &Metric{
				Name: f.Name,
				Type: f.Type,
			}
			labels []string
		)
		p2jMetric := fm.(prom2json.Metric)
		for k, v := range p2jMetric.Labels {
			labels = append(labels, fmt.Sprintf("%v=%v", k, v))
		}
		metric.Labels = strings.Join(labels, ",")
		msg, err := convertMetricToMessage(metric)
		if err != nil {
			log.Printf("error while parsing message: %v", err)
		}
		p.ch <- msg
	}
	return nil
}

func convertMetricToMessage(metric *Metric) (*Message, error) {
	b, err := json.Marshal(metric)
	if err != nil {
		return nil, fmt.Errorf("could not parse metric to JSON, %w", err)
	}

	return &Message{"", nil, b}, nil
}
