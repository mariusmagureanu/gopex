package main

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func main() {
	wait := make(chan bool)
	nc, _ := nats.Connect(nats.DefaultURL)

	fmt.Println("connected to default nats server")

	nc.Subscribe("sse", func(m *nats.Msg) {
		fmt.Printf("received a message: %s\n", string(m.Data))
	})

	<-wait
	nc.Close()
}
