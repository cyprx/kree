package kree

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/prom2json"
)

type Metric struct {
	Name   string      `json:"name"`
	Ts     string      `json:"ts"`
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
		if f.Type == dto.MetricType_HISTOGRAM.String() || f.Type == dto.MetricType_SUMMARY.String() {
			continue
		}
		err = p.parseMetric(f)
		if err != nil {
			return err
		}
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
		metric.Value = p2jMetric.Value
		metric.Ts = p2jMetric.TimestampMs
		if metric.Ts == "" {
			metric.Ts = strconv.FormatInt(time.Now().Unix(), 10)
		}
		msg, err := convertMetricToMessage(metric)
		if err != nil {
			return err
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
