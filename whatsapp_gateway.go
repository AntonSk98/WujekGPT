package main

import (
	"fmt"

	"github.com/asaskevich/EventBus"
)

// WhatsAppGateway processes received messages
type WhatsAppGateway struct {
	eventBus EventBus.Bus
}

// NewWhatsAppGateway creates a new WhatsAppGateway and subscribes handlers
func NewWhatsAppGateway(eventBus EventBus.Bus) *WhatsAppGateway {
	gateway := &WhatsAppGateway{
		eventBus: eventBus,
	}

	eventBus.Subscribe(OnMessageReceivedTopic, gateway.OnMessage)

	return gateway
}

// OnMessage handles a received message event
func (g *WhatsAppGateway) OnMessage(event OnMessageReceivedEvent) {
	fmt.Printf("[Message Received]\n")
	fmt.Printf("  ID: %s\n", event.ID)
	fmt.Printf("  Text: %q\n", event.Message)
}
