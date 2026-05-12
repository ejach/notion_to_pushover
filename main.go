package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jomei/notionapi"
)

var notionClient *notionapi.Client

type NotionWebhook struct {
	Entity struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"entity"`
}

func init() {
	token := os.Getenv("NOTION_API_KEY")
	if token == "" {
		log.Fatal("NOTION_API_KEY environment variable is required")
	}
	notionClient = notionapi.NewClient(notionapi.Token(token))
}

func getNotionPageTitle(pageID string) string {
	page, err := notionClient.Page.Get(context.Background(), notionapi.PageID(pageID))
	if err != nil {
		log.Printf("Failed to retrieve page: %v", err)
		return "Unknown Page"
	}
	for _, prop := range page.Properties {
		if titleProp, ok := prop.(*notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
			return titleProp.Title[0].PlainText
		}
	}
	return "Unknown Page"
}

// Get PUSHOVER_NOTIFICATION_TITLE from the env, else use the default
func getNotificationTitle() string {
	title := os.Getenv("PUSHOVER_NOTIFICATION_TITLE")
	if title == "" {
		return "A new page was added"
	}
	return title
}

func isTrustedRequest(c *fiber.Ctx) bool {
	body := c.Body()
	signature := c.Get("X-Notion-Signature", "")
	secret := os.Getenv("NOTION_VERIFICATION_TOKEN")
	if secret == "" {
		log.Println("NOTION_VERIFICATION_TOKEN not set")
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedSignature := "sha256=" + hex.EncodeToString(expectedMAC)

	return subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(signature)) == 1
}

func sendPushoverNotification(title, message string) error {
	form := url.Values{
		"token":   {os.Getenv("PUSHOVER_API_TOKEN")},
		"user":    {os.Getenv("PUSHOVER_USER_KEY")},
		"title":   {title},
		"message": {message},
	}

	resp, err := http.PostForm("https://api.pushover.net/1/messages.json", form)
	
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Pushover response:", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Pushover API error: %s", bodyBytes)
	}

	return nil
}

func sendNtfyNotification(title, message string) error {
	server := os.Getenv("NTFY_SERVER")
	if server == "" {
		server = "https://ntfy.sh"
	}

	topic := os.Getenv("NTFY_TOPIC")
	if topic == "" {
		return fmt.Errorf("NTFY_TOPIC environment variable is required")
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s", server, topic),
		bytes.NewBufferString(message),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Title", title)
	req.Header.Set("Content-Type", "text/plain")

	ntfyToken := os.Getenv("NTFY_BEARER_TOKEN")
	if ntfyToken != "" {
		req.Header.Set("Authorization", "Bearer "+ntfyToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("ntfy response:", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy error: %s", bodyBytes)
	}

	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	app := fiber.New(fiber.Config{
		Prefork: false,
		AppName: "Notion to Pushover",
	})

	app.Post("/", func(c *fiber.Ctx) error {
		log.Printf("Received request from %s with method %s", c.IP(), c.Method())

		if !isTrustedRequest(c) {
			log.Println("Unauthorized request")
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized request")
		}

		var data NotionWebhook
		if err := c.BodyParser(&data); err != nil {
			log.Println("Invalid payload")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid payload")
		}

		title := getNotionPageTitle(data.Entity.ID)

		notificationProvider := os.Getenv("NOTIFICATION_PROVIDER")

		switch notificationProvider {
		case "ntfy":
			if err := sendNtfyNotification(getNotificationTitle(), title); err != nil {
				log.Printf("Failed to send ntfy notification: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to send notification")
			}

		case "both":
			if err := sendPushoverNotification(getNotificationTitle(), title); err != nil {
				log.Printf("Failed to send pushover notification: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to send notification")
			}

			if err := sendNtfyNotification(getNotificationTitle(), title); err != nil {
				log.Printf("Failed to send ntfy notification: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to send notification")
			}

		default:
			if err := sendPushoverNotification(getNotificationTitle(), title); err != nil {
				log.Printf("Failed to send notification: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to send notification")
			}
		}

		log.Println("Notification sent successfully")
		return c.SendString("OK")
	})

	log.Fatal(app.Listen(":8069"))
}