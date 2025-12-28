package main

import "go.mau.fi/whatsmeow/types/events"

// WhatsAppApp defines the interface for WhatsApp application implementations
type WhatsAppApp interface {
	// Initialize sets up the WhatsApp app (authentication, QR code, etc.)
	Initialize() error

	// Shutdown gracefully closes the WhatsApp connection
	Shutdown() error

	// OnMessage handles incoming WhatsApp messages
	// Receives the entire message object, extracts data, and delegates to gateway
	OnMessage(msg *events.Message) error

	// SendMessage sends a message to a specific JID
	SendMessage(jid, text string) error
}
