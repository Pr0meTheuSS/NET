package pubsub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/lileio/pubsub"
)

var HelloTopic = "hello.topic"

type Subscriber struct{}

type HelloMsg struct {
	Greeting string
	Name     string
}

// ProtoMessage implements protoiface.MessageV1.
func (*HelloMsg) ProtoMessage() {
	panic("unimplemented")
}

// Reset implements protoiface.MessageV1.
func (*HelloMsg) Reset() {
	panic("unimplemented")
}

// String implements protoiface.MessageV1.
func (*HelloMsg) String() string {
	panic("unimplemented")
}

func PrintHello(ctx context.Context, msg *HelloMsg, m *pubsub.Msg) error {
	fmt.Printf("Message received %+v\n\n", m)

	fmt.Printf(msg.Greeting + " " + msg.Name + "\n")

	return nil
}

func (s *Subscriber) Setup(c *pubsub.Client) {
	c.On(pubsub.HandlerOptions{
		Topic:   HelloTopic,
		Name:    "print-hello",
		Handler: PrintHello,
		AutoAck: true,
		JSON:    true,
	})
}

func TestPubSub(t *testing.T) {

	pubsub.PublishJSON(context.Background(), HelloTopic, &HelloMsg{
		Greeting: "Hello",
		Name:     "Chell",
	})
	go func() {
		pubsub.Subscribe(&Subscriber{})
	}()

}
