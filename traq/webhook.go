package traq

// IWebhook traQのWebhookのinterface
type IWebhook interface {
	PostMessage(message string) error
}
