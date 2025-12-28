package main

import (
	"context"
	"fmt"
	"log"

	"github.com/asaskevich/EventBus"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// WhatsmeowApp implements WhatsAppApp using the whatsmeow library
type WhatsmeowApp struct {
	client      *whatsmeow.Client
	store       *sqlstore.Container
	dbPath      string
	authService *AuthService
	eventBus    EventBus.Bus
}

// NewWhatsmeowApp creates a new WhatsmeowApp instance with the provided gateway
func NewWhatsmeowApp(dbPath string, authFilter *AuthService, eventBus EventBus.Bus) *WhatsmeowApp {
	// Use Noop logger to minimize database logging
	container, err := sqlstore.New(context.Background(), "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbPath), waLog.Noop)
	if err != nil {
		log.Fatalf("Failed to create SQL store container: %v", err)
	}

	device, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get device: %v", err)
	}

	// Use Noop logger for client to avoid verbose logging
	client := whatsmeow.NewClient(device, waLog.Noop)

	app := &WhatsmeowApp{
		client:      client,
		store:       container,
		dbPath:      dbPath,
		authService: authFilter,
		eventBus:    eventBus,
	}

	if err := app.initialize(); err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	return app
}

// SendMessage sends a text message to the specified JID
func (w *WhatsmeowApp) SendMessage(jid, text string) error {
	parsedJID, err := types.ParseJID(jid)
	if err != nil {
		return fmt.Errorf("failed to parse JID: %w", err)
	}

	_, err = w.client.SendMessage(context.Background(), parsedJID, &waE2E.Message{
		Conversation: proto.String(text),
	})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// Shutdown gracefully closes the WhatsApp connection
func (w *WhatsmeowApp) Close() {
	if w.client != nil {
		w.client.Disconnect()
	}
	if w.store != nil {
		err := w.store.Close()
		if err != nil {
			log.Fatalf("Error closing store: %v", err)
		}
	}
}

// Initialize sets up the WhatsApp connection, handles QR code, and starts listening for messages
func (w *WhatsmeowApp) initialize() error {
	// Setup event handler before connecting
	w.setupEventHandler()

	if w.client.IsConnected() {
		return nil
	}

	// Handle QR code for login if not logged in
	if !w.client.IsLoggedIn() {
		qrChan, _ := w.client.GetQRChannel(context.Background())
		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					qr, _ := qrcode.New(evt.Code, qrcode.Medium)
					fmt.Println(qr.ToSmallString(true))
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		}()
	}

	if err := w.client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	return nil
}

// OnMessage handles incoming WhatsApp messages
func (w *WhatsmeowApp) onMessage(msg *events.Message) error {
	id := msg.Info.Chat.String()
	message := extractMessage(msg)

	// Check authentication
	if err := w.authService.CheckAuth(id, message); err != nil {
		return err
	}

	w.eventBus.Publish(OnMessageReceivedTopic, OnMessageReceivedEvent{
		ID:      id,
		Message: message,
	})

	return nil
}

// setupEventHandler registers the message handler with the client
func (w *WhatsmeowApp) setupEventHandler() {
	w.client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			if err := w.onMessage(v); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		case *events.Connected:
			fmt.Println("Connected to WhatsApp")
		case *events.LoggedOut:
			fmt.Println("Logged out from WhatsApp")
		}
	})
}

func extractMessage(msg *events.Message) string {
	msgText := ""
	if msg.Message != nil {
		if conv := msg.Message.GetConversation(); conv != "" {
			msgText = conv
		} else if extMsg := msg.Message.GetExtendedTextMessage(); extMsg != nil {
			msgText = extMsg.GetText()
		}
	}
	return msgText
}
