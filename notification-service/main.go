package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"
)

type NotificationEvent struct {
	Type      string    `json:"type"` // "due_date", "completion", etc.
	UserID    string    `json:"user_id"`
	TodoID    string    `json:"todo_id"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Connect to NATS server (message broker)
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	// Subscribe to todo events
	_, err = nc.Subscribe("todos.*", func(msg *nats.Msg) {
		var event NotificationEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println("Failed to parse event:", err)
			return
		}

		log.Printf("Received event: %+v\n", event)

		// Handle different event types
		switch event.Type {
		case "due_date":
			go sendDueDateNotification(event)
		case "completion":
			go sendCompletionNotification(event)
		default:
			log.Println("Unknown event type:", event.Type)
		}
	})
	if err != nil {
		log.Fatal("Failed to subscribe to events:", err)
	}

	// HTTP endpoints for manual testing
	e.POST("/notify/email", sendTestEmail)
	e.POST("/notify/sms", sendTestSMS)
	e.POST("/notify/push", sendTestPush)

	// Start server
	e.Logger.Fatal(e.Start(":8083"))
}

// Notification handlers
func sendDueDateNotification(event NotificationEvent) {
	message := fmt.Sprintf("Reminder: Todo '%s' is due soon", event.Title)
	sendEmail(event.UserID, "Todo Due Date Reminder", message)
	sendSMS(event.UserID, message)
	sendPushNotification(event.UserID, "Due Date Alert", message)
}

func sendCompletionNotification(event NotificationEvent) {
	message := fmt.Sprintf("Congratulations! You completed: '%s'", event.Title)
	sendEmail(event.UserID, "Todo Completed", message)
}

// Notification stubs
func sendEmail(userID, subject, body string) {
	// In production, integrate with SendGrid/Mailgun/etc.
	log.Printf("Email sent to user %s: %s - %s\n", userID, subject, body)
}

func sendSMS(userID, message string) {
	// In production, integrate with Twilio/etc.
	log.Printf("SMS sent to user %s: %s\n", userID, message)
}

func sendPushNotification(userID, title, body string) {
	// In production, integrate with FCM/APNs
	log.Printf("Push notification sent to user %s: %s - %s\n", userID, title, body)
}

// HTTP Handlers for testing
func sendTestEmail(c echo.Context) error {
	var req struct {
		UserID  string `json:"user_id"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	sendEmail(req.UserID, req.Subject, req.Body)
	return c.JSON(http.StatusOK, map[string]string{"status": "email sent"})
}

func sendTestSMS(c echo.Context) error {
	// Similar implementation to sendTestEmail
	return nil
}

func sendTestPush(c echo.Context) error {
	// Similar implementation to sendTestEmail
	return nil
}
