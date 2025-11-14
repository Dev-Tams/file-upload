package storage

// import (
// 	"mime/multipart"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/s3/s3manager"
// 	"github.com/dev-tams/file-upload/internal/config"
// )

// type S3Storage struct {
// 	bucket string
// }

// func NewS3Storage(cfg *config.AppConfig) (StorageProvider, error) {
// 	return &S3Storage{
// 		bucket: cfg.S3Bucket,
// 	}, nil
// }

// func (s *S3Storage) Save(file *multipart.FileHeader, storedName string) (string, error) {
// 	f, err := file.Open()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer f.Close()

// 	_, err = s.Uploader.Upload(&s3manager.UploadInput{
// 		Bucket: aws.String(s.bucket),
// 		Key:    aws.String(storedName),
// 		Body:   f,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	return storedName, nil
// }

// func (s *S3Storage) Delete(path string) error {
// 	return nil // implement when needed
// }
