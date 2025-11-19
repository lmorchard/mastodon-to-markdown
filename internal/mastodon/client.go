package mastodon

import (
	"context"
	"fmt"

	"github.com/lmorchard/mastodon-to-markdown/internal/config"
	"github.com/mattn/go-mastodon"
)

// Client wraps the go-mastodon client with our configuration
type Client struct {
	client *mastodon.Client
	config *config.Config
}

// NewClient creates a new Mastodon API client
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.Mastodon.Server == "" {
		return nil, fmt.Errorf("mastodon server URL is required")
	}
	if cfg.Mastodon.AccessToken == "" {
		return nil, fmt.Errorf("mastodon access token is required")
	}

	client := mastodon.NewClient(&mastodon.Config{
		Server:      cfg.Mastodon.Server,
		AccessToken: cfg.Mastodon.AccessToken,
	})

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

// VerifyCredentials checks if the access token is valid and returns account info
func (c *Client) VerifyCredentials(ctx context.Context) (*mastodon.Account, error) {
	account, err := c.client.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify credentials: %w", err)
	}
	return account, nil
}

// GetStatuses fetches statuses for the authenticated user's account
// maxID and sinceID are used for pagination
func (c *Client) GetStatuses(ctx context.Context, accountID mastodon.ID, pg *mastodon.Pagination) ([]*mastodon.Status, error) {
	statuses, err := c.client.GetAccountStatuses(ctx, accountID, pg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch statuses: %w", err)
	}
	return statuses, nil
}

// GetFavourites fetches favorited posts for the authenticated user
func (c *Client) GetFavourites(ctx context.Context, pg *mastodon.Pagination) ([]*mastodon.Status, error) {
	statuses, err := c.client.GetFavourites(ctx, pg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch favourites: %w", err)
	}
	return statuses, nil
}

// GetClient returns the underlying go-mastodon client for advanced use
func (c *Client) GetClient() *mastodon.Client {
	return c.client
}
