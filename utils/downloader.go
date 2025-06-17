package utils

import (
	"context"
	s3client "dhcli/configs"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadHTTPFile function for get a file from http or https
func DownloadHTTPFile(url string, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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

	if err := s3Client.DownloadFile(ctx, bucket, key, localPath); err != nil {
		return fmt.Errorf("S3 download failed: %w", err)
	}

	return nil
}
