package utils

import (
	"context"
	s3client "dhcli/configs"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadHTTPFile function for get a file from http or https
func DownloadHTTPFile(url string, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// DownloadS3File Function for download file or directory form S3
func DownloadS3File(s3Client *s3client.Client, ctx context.Context,
	parsedPath *ParsedPath, localPath string,
) error {
	bucket := parsedPath.Host
	key := parsedPath.Path

	// Se è una "cartella"
	if strings.HasSuffix(key, "/") {
		files, err := s3Client.ListFiles(ctx, bucket, key, aws.Int32(200))
		if err != nil {
			return fmt.Errorf("failed to list S3 folder: %w", err)
		}

		for _, file := range files {
			// Costruisci path locale relativo
			relativePath := strings.TrimPrefix(file.Path, key)
			targetPath := filepath.Join(localPath, relativePath)

			// Crea directory se necessario
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create local directory: %w", err)
			}

			if err := s3Client.DownloadFile(ctx, bucket, file.Path, targetPath); err != nil {
				return fmt.Errorf("failed to download file: %w", err)
			}
		}
	} else {
		// File singolo
		if err := s3Client.DownloadFile(ctx, bucket, key, localPath); err != nil {
			return fmt.Errorf("S3 download failed: %w", err)
		}
	}

	return nil
}
