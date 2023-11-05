package util

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

func Upload(ctx context.Context, s3Client *s3.Client, input *s3.PutObjectInput) {
	_, err := s3Client.PutObject(ctx, input)
	if err != nil {
		log.Fatal("put object failed.", err)
	}
	log.Println("upload completed")
}

func GetFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("open file failed.", err)
	}

	return file
}
