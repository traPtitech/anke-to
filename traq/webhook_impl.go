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

	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
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

	sb := &strings.Builder{}
	_, err = io.Copy(sb, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Message sent to %s, message: %s, response: %s\n", url, message, sb.String())

	return nil
}

func calcHMACSHA1(message string) string {
	mac := hmac.New(sha1.New, []byte(os.Getenv("TRAQ_WEBHOOK_SECRET")))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
