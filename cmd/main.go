package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/esuEdu/aws-s3-go/storage"
)

func main() {
	cmd := flag.String("cmd", "", "Command: upload, list, download")
	bucket := flag.String("bucket", "", "s3 bucket name")
	file := flag.String("file", "", "Path to local file (upload/download)")
	key := flag.String("key", "", "S3 object key (upload/download)")

	flag.Parse()

	if *cmd == "" || *bucket == "" {
		log.Fatal("Required flags: -cmd and -bucket")
	}

	ctx := context.Background()
	client, err := storage.NewS3Client(ctx, *bucket)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	switch *cmd {
	case "upload":
		if *file == "" || *key == "" {
			log.Fatal("Upload requires -file and -key")
		}
		if err := client.UploadFile(ctx, *key, *file); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}
		fmt.Println("Upload successful")

	case "list":
		if err := client.ListObjects(ctx); err != nil {
			log.Fatalf("List failed: %v", err)
		}

	case "download":
		if *file == "" || *key == "" {
			log.Fatal("Download requires -file and -key")
		}
		if err := client.DownloadFile(ctx, *key, *file); err != nil {
			log.Fatalf("Download failed: %v", err)
		}
		fmt.Println("Download successful")

	default:
		log.Fatalf("Unknown command: %s", *cmd)
	}
}
