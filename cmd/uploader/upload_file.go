package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func uploadFile(fileName string, max <-chan struct{}, wg *sync.WaitGroup, filesErrs chan<- string) {
	defer wg.Done()
	defer func() { <-max }()

	fullName := fmt.Sprintf("./temp/%s", fileName)
	f, err := os.Open(fullName)
	if err != nil {
		log.Printf("Error: open file %s\n", fileName)
		return
	}
	defer f.Close()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fullName),
		Body:   f,
	})
	if err != nil {
		log.Printf("Error: upload file %s\n", fileName)
		filesErrs <- fileName
		return
	}

	fmt.Printf("file: %s uploaded!\n", fullName)
}
