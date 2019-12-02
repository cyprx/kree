package main

import (
	"os"
	"os/signal"

	"github.com/cyprx/kree"
)

var (
	KAFKA_BROKERS_ADDR = []string{"", ""}
	KAFKA_METRIC_TOPIC = "data.prometheus.metrics"
)

func main() {

	repo := kree.NewEndpointRepository()
	ch := make(chan *kree.Message, 1)
	q := make(chan os.Signal, 1)
	signal.Notify(q, os.Interrupt)

	collector := kree.NewCollector(repo, ch, q)
	publisher := kree.NewPublisher(KAFKA_METRIC_TOPIC, KAFKA_BROKERS_ADDR, ch, q)
	go collector.Run()
	go publisher.Publish()

	<-q
}
