package s3client

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3 *s3.Client
}

type Config struct {
	AccessKey   string
	SecretKey   string
	AccessToken string
	Region      string
	EndpointURL string
}

func NewClient(ctx context.Context, cfgCreds Config) (*Client, error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		cfgCreds.AccessKey,
		cfgCreds.SecretKey,
		cfgCreds.AccessToken,
	))

	// Load AWS configuration with credentials and region
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(cfgCreds.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Initialize S3 client options
	s3Options := func(o *s3.Options) {
		// If a custom endpoint is provided, set it as BaseEndpoint
		if cfgCreds.EndpointURL != "" {
			o.BaseEndpoint = aws.String(cfgCreds.EndpointURL)
			o.UsePathStyle = true // Necessary for some S3-compatible services
		}
	}

	// Create S3 client with the specified options
	return &Client{
		s3: s3.NewFromConfig(cfg, s3Options),
	}, nil
}

func ParseS3Path(s3path string) (bucket string, key string, filename string, err error) {
	const prefix = "s3://"
	if !strings.HasPrefix(s3path, prefix) {
		return "", "", "", fmt.Errorf("invalid s3 path: must start with %q", prefix)
	}

	trimmed := strings.TrimPrefix(s3path, prefix)
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid s3 path: must contain bucket and key")
	}

	bucket = parts[0]
	key = parts[1]

	// Extract file name from the key
	keyParts := strings.Split(key, "/")
	filename = keyParts[len(keyParts)-1]

	return bucket, key, filename, nil
}

// DownloadFile downloads a file from S3 and saves it locally
func (c *Client) DownloadFile(ctx context.Context, bucket, key, localPath string) error {

	fmt.Printf("Downloading from S3 path: s3://%s/%s\n", bucket, key)
	fmt.Printf("Saving to local path: %s\n", localPath)

	output, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer output.Body.Close()

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, output.Body)
	if err != nil {
		return fmt.Errorf("failed to write to local file: %w", err)
	}

	return nil
}

//// UploadFile uploads a local file to the specified S3 bucket and key
//func (c *Client) UploadFile(ctx context.Context, bucket, key, localPath string) error {
//	// Open the local file for reading
//	file, err := os.Open(localPath)
//	if err != nil {
//		return fmt.Errorf("failed to open local file: %w", err)
//	}
//	defer file.Close()
//
//	// Get file info for content length (optional but recommended)
//	fileInfo, err := file.Stat()
//	if err != nil {
//		return fmt.Errorf("failed to stat local file: %w", err)
//	}
//
//	// Upload the file to S3
//	_, err = c.s3.PutObject(ctx, &s3.PutObjectInput{
//		Bucket:        &bucket,
//		Key:           &key,
//		Body:          file,
//		ContentLength: fileInfo.Size(),
//	})
//	if err != nil {
//		return fmt.Errorf("failed to upload file to S3: %w", err)
//	}
//
//	return nil
//}
