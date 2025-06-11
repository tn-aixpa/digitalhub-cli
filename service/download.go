// service/download.go
package service

import (
	"context"
	s3client "dhcli/configs"
	"dhcli/models"
	"dhcli/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func DownloadFileWithOptions(env, output, project, name, resource, id string, originalArgs []string) error {
	validResource := utils.TranslateEndpoint(resource)

	if validResource != "projects" && project == "" {
		return errors.New("project is mandatory when performing this operation on resources other than projects")
	}

	params := map[string]string{}
	if id == "" {
		if name == "" {
			return errors.New("you must specify id or name")
		}
		params["name"] = name
		params["versions"] = "latest"
	}

	_, section := utils.LoadIniConfig([]string{env})
	method := "GET"
	url := utils.BuildCoreUrl(section, project, validResource, id, params)
	req := utils.PrepareRequest(method, url, nil, section.Key("access_token").String())
	body, err := utils.DoRequest(req)

	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var resp models.Response[models.Artifact] //TODO fix the type in future

	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	if len(resp.Content) == 0 {
		return fmt.Errorf("no artifact was found in Content response")
	}

	ctx := context.Background()
	cfg := s3client.Config{
		AccessKey:   section.Key("aws_access_key_id").String(),
		SecretKey:   section.Key("aws_secret_access_key").String(),
		AccessToken: section.Key("aws_session_token").String(),
		Region:      section.Key("aws_region").String(),
		EndpointURL: section.Key("aws_endpoint_url").String(),
	}
	client, err := s3client.NewClient(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	for i, artifact := range resp.Content {
		fmt.Printf("Artifact #%d - Path: %s\n", i+1, artifact.Spec.Path)

		bucket, key, filename, err := s3client.ParseS3Path(artifact.Spec.Path)
		fmt.Println(bucket, key, filename)
		if err != nil {
			return fmt.Errorf("invalid S3 path: %w", err)
		}

		localPath := filename
		if err := client.DownloadFile(ctx, bucket, key, localPath); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
	}

	log.Println("File downloaded successfully.")
	return nil
}
