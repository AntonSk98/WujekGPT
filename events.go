package main

// Event topics
const OnMessageReceivedTopic = "message.received"
const OnResponseReadyTopic = "response.ready"

// OnMessageReceivedEvent represents an event when a message is received
type OnMessageReceivedEvent struct {
	ID      string
	Message string
}
