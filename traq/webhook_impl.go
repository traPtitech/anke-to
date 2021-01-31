package traq

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"

	"fmt"
	"net/http"
	netUrl "net/url"
	"os"
	"strings"

	"github.com/labstack/echo"
)

// Webhook Webhookの構造体
type Webhook struct{}

// PostMessage Webhookでのメッセージの投稿
func (*Webhook) PostMessage(message string) error {
	url := "https://q.trap.jp/api/v3/webhooks/" + os.Getenv("TRAQ_WEBHOOK_ID")
	req, err := http.NewRequest("POST",
		url,
		strings.NewReader(message))
	if err != nil {
		return err
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	req.Header.Set("X-TRAQ-Signature", calcHMACSHA1(message))

	query := netUrl.Values{}
	query.Add("embed", "1")
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	response := make([]byte, 512)
	_, err = resp.Body.Read(response)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Message sent to %s, message: %s, response: %s\n", url, message, response)

	return nil
}

func calcHMACSHA1(message string) string {
	mac := hmac.New(sha1.New, []byte(os.Getenv("TRAQ_WEBHOOK_SECRET")))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
