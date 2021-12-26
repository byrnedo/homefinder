package repos

import (
	"bufio"
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Void struct{}

type HistoryRepo interface {
	GetHistory(ctx context.Context) (map[string]Void, error)
	SaveHistory(ctx context.Context, keys map[string]Void) error
}

type S3HistoryRepo struct {
	s3c    *s3.Client
	bucket string
}

func NewS3HistoryRepo(s3c *s3.Client, bucket string) S3HistoryRepo {
	return S3HistoryRepo{
		s3c:    s3c,
		bucket: bucket,
	}
}

func (s S3HistoryRepo) GetHistory(ctx context.Context) (map[string]Void, error) {
	obj, err := s.s3c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String("listings"),
	})
	if err != nil {
		return nil, err
	}

	defer obj.Body.Close()

	m := map[string]Void{}
	scanner := bufio.NewScanner(obj.Body)
	for scanner.Scan() {
		text := scanner.Text()
		m[text] = Void{}
	}
	return m, nil
}

func (s S3HistoryRepo) SaveHistory(ctx context.Context, keys map[string]Void) error {

	writer := &bytes.Buffer{}

	for k := range keys {
		_, err := writer.WriteString(k + "\n")
		if err != nil {
			return err
		}
	}

	log.Println("saving to s3:")
	_, err := s.s3c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String("listings"),
		Body:   writer,
	})
	return err
}

type EmptyHistoryRepo struct {
}

func (e EmptyHistoryRepo) GetHistory(ctx context.Context) (map[string]Void, error) {
	return map[string]Void{}, nil
}

func (e EmptyHistoryRepo) SaveHistory(ctx context.Context, keys map[string]Void) error {
	return nil
}
