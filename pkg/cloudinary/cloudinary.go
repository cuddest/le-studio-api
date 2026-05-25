package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Client provides cloudinary upload methods.
type Client struct {
	cld *cloudinary.Cloudinary
}

// New creates a cloudinary helper.
func New(cloudName, apiKey, apiSecret string) (*Client, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}
	return &Client{cld: cld}, nil
}

// UploadURL uploads from URL string and returns secure URL.
func (c *Client) UploadURL(ctx context.Context, source string) (string, error) {
	result, err := c.cld.Upload.Upload(ctx, source, uploader.UploadParams{})
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}
	return result.SecureURL, nil
}

// UploadFile uploads a file and returns secure URL.
func (c *Client) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	result, err := c.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: "le-studio/coaches",
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}
	return result.SecureURL, nil
}
