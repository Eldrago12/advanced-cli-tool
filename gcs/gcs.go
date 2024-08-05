package gcs

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"google.golang.org/api/option"
)

func UploadFilesToBucket(bucketName string) {
	viper.SetConfigFile(".gcp_config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	keyFilePath := viper.GetString("gcp_key_file")
	if keyFilePath == "" {
		log.Fatalf("GCP key file not found in config")
	}

	projectID := viper.GetString("gcp_project_id")
	if projectID == "" {
		log.Fatalf("GCP project ID not found in config")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(keyFilePath))
	if err != nil {
		log.Fatalf("Failed to create GCP storage client: %v", err)
	}

	bucket := client.Bucket(bucketName)
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	files, err := os.ReadDir(pwd)
	if err != nil {
		log.Fatalf("Failed to read current directory: %v", err)
	}

	p := mpb.New(mpb.WithWidth(64))
	totalFiles := len(files)
	bar := p.AddBar(int64(totalFiles),
		mpb.PrependDecorators(
			decor.Name("Uploading files: "),
			decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncWidth),
		),
	)

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			filePath := filepath.Join(pwd, fileName)

			go func() {
				defer bar.Increment()
				f, err := os.Open(filePath)
				if err != nil {
					log.Printf("Failed to open file %s: %v", fileName, err)
					return
				}
				defer f.Close()

				ctx, cancel := context.WithTimeout(ctx, time.Second*50)
				defer cancel()

				wc := bucket.Object(fileName).NewWriter(ctx)
				if _, err := wc.Write(f); err != nil {
					log.Printf("Failed to write file %s: %v", fileName, err)
					return
				}
				if err := wc.Close(); err != nil {
					log.Printf("Failed to close writer for file %s: %v", fileName, err)
					return
				}

				fmt.Printf("Uploaded file: %s\n", fileName)
			}()
		}
	}

	p.Wait()
}
