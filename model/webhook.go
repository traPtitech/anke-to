package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"

	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
)

func CalcHMACSHA1(message string) string {
	mac := hmac.New(sha1.New, []byte(os.Getenv("TRAQ_WEBHOOK_SECRET")))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func PostMessage(c echo.Context, message string) error {
	url := "https://q.trap.jp/api/1.0/webhooks/" + os.Getenv("TRAQ_WEBHOOK_ID")
	req, err := http.NewRequest("POST",
		url,
		strings.NewReader(message))
	if err != nil {
		return err
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	req.Header.Set("X-TRAQ-Signature", CalcHMACSHA1(message))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	response := make([]byte, 512)
	resp.Body.Read(response)

	fmt.Printf("Message sent to %s, message: %s, response: %s\n", url, message, response)

	return c.NoContent(http.StatusNoContent)
}
