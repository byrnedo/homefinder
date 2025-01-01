package repos

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Void struct{}

type HistoryRepo interface {
	GetHistory(ctx context.Context) ([]agents.Listing, error)
	SaveHistory(ctx context.Context, listings []agents.Listing) error
}

type S3HomesBjelinHistoryRepo struct {
	s3c    *s3.Client
	bucket string
}

type S3HomesHistoryRepo struct {
	s3c    *s3.Client
	bucket string
}

func NewS3HomesHistoryRepo(s3c *s3.Client, bucket string) S3HomesHistoryRepo {
	return S3HomesHistoryRepo{
		s3c:    s3c,
		bucket: bucket,
	}
}

func (s S3HomesHistoryRepo) GetHistory(ctx context.Context) (map[string]Void, error) {
	return get(ctx, s.s3c, s.bucket, "listings")
}

func (s S3HomesHistoryRepo) SaveHistory(ctx context.Context, keys map[string]Void) error {
	return save(ctx, s.s3c, s.bucket, "listings", keys)
}

type S3JobsHistoryRepo struct {
	s3c    *s3.Client
	bucket string
}

func get(ctx context.Context, s3c *s3.Client, bucket, filename string) (map[string]Void, error) {

	obj, err := s3c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    aws.String(filename),
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

func save(ctx context.Context, s3c *s3.Client, bucket, filename string, keys map[string]Void) error {

	orig, err := get(ctx, s3c, bucket, filename)
	if err != nil {
		return err
	}

	writer := &bytes.Buffer{}

	for k, v := range keys {
		orig[k] = v
	}

	for k := range orig {
		_, err := writer.WriteString(k + "\n")
		if err != nil {
			return err
		}
	}

	_, err = s3c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    aws.String(filename),
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

type FileHistoryRepo struct {
	Name string
}

func (e FileHistoryRepo) GetHistory(ctx context.Context) (list []agents.Listing, err error) {
	b, _ := os.ReadFile(e.Name)

	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	for scanner.Scan() {
		text := scanner.Text()
		l := agents.Listing{}
		err := json.Unmarshal([]byte(text), &l)
		if err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func (e FileHistoryRepo) SaveHistory(ctx context.Context, list []agents.Listing) error {
	f, err := os.OpenFile(e.Name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_ = f.Truncate(0)
	_, _ = f.Seek(0, 0)

	for _, listing := range list {
		b, err := json.Marshal(listing)
		if err != nil {
			return err
		}

		_, err = f.WriteString(string(b) + "\n")
		if err != nil {
			return err
		}
	}
	return nil

}
