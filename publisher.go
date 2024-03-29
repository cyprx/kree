package kree

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
)

// Publisher will listen to message channel and publish any messages it received
type Publisher interface {
	Publish()
}

type publisher struct {
	topic string
	ap    sarama.AsyncProducer
	input chan *Message
	quit  chan os.Signal
}

// NewPublisher return publisher implementation object
func NewPublisher(topic string, brokers []string, i chan *Message, q chan os.Signal) Publisher {
	ap, err := sarama.NewAsyncProducer(brokers, nil)
	if err != nil {
		log.Fatalf("could not start async producer, err: %v", err)
	}

	return &publisher{topic, ap, i, q}
}

// Publish start listening to input channel for new messages
func (p *publisher) Publish() {
	log.Println("Starting publisher...")
	for {
		select {
		case m := <-p.input:
			pm := &sarama.ProducerMessage{
				Topic: p.topic,
				Key:   sarama.ByteEncoder(m.Key),
				Value: sarama.ByteEncoder(m.Value),
			}
			p.ap.Input() <- pm
		case err := <-p.ap.Errors():
			log.Printf("failed to produce message %v", err)
		case <-p.quit:
			log.Println("Receiving stop signal, closing Publisher")
			break
		}
	}
}
