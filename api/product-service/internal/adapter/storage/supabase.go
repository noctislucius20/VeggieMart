package storage

import (
	"io"
	"product-service/config"

	"github.com/labstack/gommon/log"
	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseInterface interface {
	UploadFile(path string, file io.Reader) (string, error)
}

type supabaseStruct struct {
	cfg    *config.Config
	logger *log.Logger
}

// UploadFile implements SupabaseInterface.
func (s *supabaseStruct) UploadFile(path string, file io.Reader) (string, error) {
	client := storage_go.NewClient(s.cfg.Storage.Url, s.cfg.Storage.Key, map[string]string{"Content-Type": "image/png"})

	_, err := client.UploadFile(s.cfg.Storage.Bucket, path, file)
	if err != nil {
		s.logger.Errorf("failed to upload file: %v", err)
		return "", err
	}

	result := client.GetPublicUrl(s.cfg.Storage.Bucket, path)

	return result.SignedURL, nil
}

func NewSupabase(cfg *config.Config, logger *log.Logger) SupabaseInterface {
	return &supabaseStruct{
		cfg:    cfg,
		logger: logger,
	}
}
