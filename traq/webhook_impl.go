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

// traqMessageLimit traQのメッセージの最大文字数
const traqMessageLimit = 2000

// PostMessage Webhookでのメッセージの投稿
func (*Webhook) PostMessage(message string) error {
	for _, chunk := range splitMessage(message, traqMessageLimit) {
		if err := postSingleMessage(chunk); err != nil {
			return err
		}
	}
	return nil
}

func postSingleMessage(message string) error {
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

// splitMessage はメッセージをlimit文字以内のチャンクに分割する。
// まず改行・スペース単位でトークン化し、各トークンを貪欲にチャンクへ詰める。
// 改行はトークン間のセパレーターとして保持されるため、元の構造をできるだけ維持する。
func splitMessage(message string, limit int) []string {
	if len([]rune(message)) <= limit {
		return []string{message}
	}

	// トークンはテキストと直前のセパレーター（""、"\n"、" "）のペア
	type token struct {
		text string
		sep  string
	}

	var tokens []token
	lines := strings.Split(message, "\n")
	for i, line := range lines {
		words := strings.Split(line, " ")
		for j, word := range words {
			sep := " "
			if j == 0 && i == 0 {
				sep = ""
			} else if j == 0 {
				sep = "\n"
			}
			tokens = append(tokens, token{text: word, sep: sep})
		}
	}

	var chunks []string
	current := &strings.Builder{}

	for _, tok := range tokens {
		sep := tok.sep
		if current.Len() == 0 {
			sep = ""
		}
		potentialLen := len([]rune(current.String())) + len([]rune(sep)) + len([]rune(tok.text))
		if current.Len() == 0 || potentialLen <= limit {
			current.WriteString(sep)
			current.WriteString(tok.text)
		} else {
			chunks = append(chunks, current.String())
			current.Reset()
			current.WriteString(tok.text)
		}
	}

	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks
}

func calcHMACSHA1(message string) (string, error) {
	mac := hmac.New(sha1.New, []byte(os.Getenv("TRAQ_WEBHOOK_SECRET")))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to write message to mac: %w", err)
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}
