package kree

import (
	"context"
	"testing"
)

func TestMakeRequest(t *testing.T) {
	c := NewCollector()
	ctx := context.Background()

	e := &endpoint{domain: "http://gum-delivery.tiki.services", path: "/metrics"}

	c.makeRequest(ctx, e)
}
