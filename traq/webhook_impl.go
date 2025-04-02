package traq

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"

	"fmt"
	"net/http"
	netUrl "net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// Webhook Webhookの構造体
type Webhook struct{}

// NewWebhook Webhookのコンストラクター
func NewWebhook() *Webhook {
	return new(Webhook)
}

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
	messageHMAC, err := calcHMACSHA1(message)
	if err != nil {
		return err
	}
	req.Header.Set("X-TRAQ-Signature", messageHMAC)

	query := netUrl.Values{}
	query.Add("embed", "1")
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sb := &strings.Builder{}
	_, err = io.Copy(sb, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Message sent to %s, message: %s, response: %s\n", url, message, sb.String())

	return nil
}

func calcHMACSHA1(message string) (string, error) {
	mac := hmac.New(sha1.New, []byte(os.Getenv("TRAQ_WEBHOOK_SECRET")))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to write message to mac: %w", err)
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}
