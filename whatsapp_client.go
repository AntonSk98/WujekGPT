package main

// SendMessageCommand represents a command to send a message
type SendMessageCommand struct {
	ID   string
	Text string
}

// WhatsAppClient wraps the WhatsAppApp to provide a clean interface for sending messages
type WhatsAppClient struct {
	app WhatsAppApp
}

// NewWhatsAppClient creates a new WhatsAppClient with the provided app
func NewWhatsAppClient(app WhatsAppApp) *WhatsAppClient {
	return &WhatsAppClient{
		app: app,
	}
}

// SendMessage sends a message via the WhatsAppApp
func (w *WhatsAppClient) SendMessage(command *SendMessageCommand) error {
	return w.app.SendMessage(command.ID, command.Text)
}
