package main

import (
	"bytes"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/donalbyrne/homefinder/internal/app"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	bucket := os.Getenv("BUCKET")
	lambda.Start(func(ctx context.Context) {
		// download prev from s3
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("fetching file:")
		s3c := s3.NewFromConfig(cfg)
		obj, err := s3c.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("listings"),
		})
		if err != nil {
			log.Println("WARNING: error getting file from s3:", err)
		} else {
			log.Println("writing local file:")
			b, _ := ioutil.ReadAll(obj.Body)
			if err := ioutil.WriteFile(app.FileName, b, 0644); err != nil {
				log.Fatal(err)
			}
			log.Println("writing local file: done")
		}
		log.Println("fetching file: done")
		// put in file as listings-seen

		// run app
		log.Println("running")
		app.Run(ctx)

		b, err := ioutil.ReadFile(app.FileName)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("saving to s3:")
		_, err = s3c.PutObject(ctx, &s3.PutObjectInput{
			Bucket: &bucket,
			Key:    aws.String("listings"),
			Body:   bytes.NewReader(b),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("saving to s3: done")
	})
}
