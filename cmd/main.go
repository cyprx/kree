package main

import (
	"os"
	"os/signal"
	"strings"

	"github.com/cyprx/kree"
)

var (
	KAFKA_BROKERS_ADDR   = strings.Split(getEnv("KAFKA_BROKERS_ADDR", "localhost:2081"), ",")
	KAFKA_METRIC_TOPIC   = getEnv("KAFKA_METRIC_TOPIC", "topic_test")
	MYSQL_CONNECTION_URI = getEnv("MYSQL_CONNECTION_URI", "root:password@tcp(localhost:3306)/test")
)

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		val = fallback
	}
	return val
}

func main() {
	kree.InitDB(MYSQL_CONNECTION_URI)
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
