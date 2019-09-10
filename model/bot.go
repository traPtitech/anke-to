package model

import (
	"net/http"
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"
)


// 投稿先のチャンネルのID
const channelId = "76da7291-eec8-411d-9ed2-8a221821835f"


func PostMessageUsingBot(c echo.Context, message string) error {
	url := "https://q.trap.jp/api/1.0/channels/" + channelId + "/messages"
	req, err := http.NewRequest("POST",
		url,
		strings.NewReader(message)
	)
	if err != nil {
		return err
	}

	req.Header.Set(echo.HeaderAuthorization, "Bearer " + os.Getenv("TRAQ_BOT_ACCESSTOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	response := make([]byte, 512)
	resp.Body.Read(response)

	fmt.Printf("Message sent to %s, message: %s, response: %s\n", url, message, response)

	return nil
}