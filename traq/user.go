//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package traq

import "golang.org/x/oauth2"

type IUser interface {
	GetMyID(token *oauth2.Token) (string, error)
}
