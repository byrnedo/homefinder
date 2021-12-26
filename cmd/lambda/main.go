package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/byrnedo/homefinder/internal/app"
	"github.com/byrnedo/homefinder/internal/pkg/repos"
)

func main() {
	bucket := os.Getenv("BUCKET")
	lambda.Start(func(ctx context.Context) {
		// download prev from s3
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatal(err)
		}

		s3c := s3.NewFromConfig(cfg)

		historyRepo := repos.NewS3HistoryRepo(s3c, bucket)

		app.Run(ctx, historyRepo)
	})
}
