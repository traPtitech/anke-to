//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package traq

// IWebhook traQのWebhookのinterface
type IWebhook interface {
	PostMessage(message string) error
}
