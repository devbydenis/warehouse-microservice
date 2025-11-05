package storage

import (
	"context"
	"fmt"
	"micro-warehouse/warehouse-service/configs"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseInterface interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error)
}

type SupabaseStorage struct {
	client *storage_go.Client
	cfg    configs.Config
}

// UploadFile implements SupabaseInterface.
func (s *SupabaseStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate uinique filename
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, ext), timestamp, ext)

	// Create file path
	filePath := fmt.Sprintf("%s/%s", folder, filename)

	// Use the simpler implementation with proper Content-Type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		// Set default content type based on file extension
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		case ".svg":
			contentType = "image/svg+xml"
		default:
			contentType = "application/octet-stream"
		}
	}

	// Create client with proper Content-Type
	client := storage_go.NewClient(s.cfg.Supabase.Url, s.cfg.Supabase.Key, map[string]string{
		"Content-Type": contentType,
	})

	// Upload file
	_, err = client.UploadFile(s.cfg.Supabase.Bucket, filePath, src)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to supabase: %w", err)
	}

	// Get public URL
	publicUrl := client.GetPublicUrl(s.cfg.Supabase.Bucket, filePath)

	return &UploadResult{
		URL:      publicUrl.SignedURL,
		Path:     filePath,
		Filename: filename,
	}, nil

}

type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

func NewSupabaseStorage(cfg configs.Config) SupabaseInterface {
	client := storage_go.NewClient(cfg.Supabase.Url, cfg.Supabase.Key, nil)
	return &SupabaseStorage{
		client: client,
		cfg:    cfg,
	}
}
