package server

import (
	"os"
	"testing"
	"time"

	"github.com/porthos-rpc/porthos-go/broker"
)

func TestMetricsShipperExtension(t *testing.T) {
	b, _ := broker.NewBroker(os.Getenv("AMQP_URL"))

	ext := NewMetricsShipperExtension(b, MetricsShipperConfig{BufferSize: 2})

	ext.outgoing <- OutgoingRPC{&request{serviceName: "SampleService", methodName: "test1"}, &response{}, 4 * time.Millisecond, 200}
	ext.outgoing <- OutgoingRPC{&request{serviceName: "SampleService", methodName: "test2"}, &response{}, 5 * time.Millisecond, 201}
	ext.outgoing <- OutgoingRPC{&request{serviceName: "SampleService", methodName: "test2"}, &response{}, 6 * time.Millisecond, 201}
	ext.outgoing <- OutgoingRPC{&request{serviceName: "SampleService", methodName: "test3"}, &response{}, 7 * time.Millisecond, 202}
	ext.outgoing <- OutgoingRPC{&request{serviceName: "SampleService", methodName: "test4"}, &response{}, 8 * time.Millisecond, 200}
	ext.outgoing <- OutgoingRPC{&request{serviceName: "SampleService", methodName: "test4"}, &response{}, 9 * time.Millisecond, 200}

	ch, _ := b.OpenChannel()

	dc, _ := ch.Consume(
		"porthos.metrics", // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)

	shippedMetricsCount := 0

	go func() {
		for range dc {
			shippedMetricsCount++
		}
	}()

	<-time.After(2 * time.Second)

	if shippedMetricsCount != 3 {
		t.Errorf("Excepted 3 shipped metrics, got %d", shippedMetricsCount)
	}
}
