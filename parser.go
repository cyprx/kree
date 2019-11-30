package kree

import (
	"encoding/json"
	"io"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/prom2json"
)

type parser struct {
	p expfmt.TextParser
}

func newParser() *parser {
	return &parser{expfmt.TextParser{}}
}

func (p *parser) parseTextToJSON(line io.Reader) ([]byte, error) {
	result := []*prom2json.Family{}
	mfs, err := p.p.TextToMetricFamilies(line)
	if err != nil {
		return []byte{}, err
	}
	for _, mf := range mfs {
		result = append(result, prom2json.NewFamily(mf))
	}

	return json.Marshal(result)
}
