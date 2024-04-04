package client

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/script/v1"
)

// NewWithAccessToken creates the Apps Script client
func NewWithAccessToken(ctx context.Context, tok string) (*script.Service, error) {
	oauthCfg := &oauth2.Config{}
	oauthToken := &oauth2.Token{
		AccessToken: tok,
	}
	client := oauthCfg.Client(ctx, oauthToken)

	// Create Apps Script service
	return script.NewService(ctx, option.WithHTTPClient(client))
}
