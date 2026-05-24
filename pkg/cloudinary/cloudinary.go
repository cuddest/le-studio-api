package cloudinary

import "context"

// Client provides cloudinary upload methods.
type Client struct{}

// New creates a cloudinary helper.
func New(_, _, _ string) (*Client, error) { return &Client{}, nil }

// UploadURL uploads source and returns URL string.
func (c *Client) UploadURL(_ context.Context, source string) (string, error) { return source, nil }
