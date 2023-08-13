package main

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var (
	s3Client *s3.S3
	s3Bucket string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: load .env")
	}

	s, err := session.NewSession(
		&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("ID"),
				os.Getenv("SECRET"),
				"",
			),
		},
	)

	if err != nil {
		log.Fatal("Error: session aws")
	}

	s3Client = s3.New(s)
	s3Bucket = "simple-upload-bucket"
}

func main() {
	max := make(chan struct{}, 10)
	filesErrs := make(chan string, 3)

	d, err := os.Open("./temp")
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		close(max)
		close(filesErrs)
		defer d.Close()
	}()

	wg := &sync.WaitGroup{}

	go func() {
		for name := range filesErrs {
			max <- struct{}{}
			wg.Add(1)
			go uploadFile(name, max, wg, filesErrs)
		}
	}()

	for {
		files, err := d.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error: read dir \n")
			continue
		}
		max <- struct{}{}
		wg.Add(1)
		go uploadFile(files[0].Name(), max, wg, filesErrs)
	}

	wg.Wait()
}
