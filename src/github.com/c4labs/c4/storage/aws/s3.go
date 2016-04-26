package aws

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (r *AWSResource) inS3(key string) bool {
	resp, err := r.srv.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String("test.cccc.io"), // Required
		Key:    aws.String(key),            // Required
	})
	if resp.ContentLength == nil {
		return false
	}
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (r *AWSResource) upload_file(path string) {

	if r.inS3(filepath.Base(path)) {
		return
	}
	fmt.Printf("upload_file: %s\n", filepath.Base(path))
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open file", err)
	}
	_, err = r.uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String("test.cccc.io"),
		Key:    aws.String(filepath.Base(path)),
	})
	if err != nil {
		fmt.Println("Failed to upload", err)
	}
}

func file_downloader(c chan SrcDest, next chan SrcDest) error {
	downloader := s3manager.NewDownloader(session.New())
	for {
		sd := <-c
		fmt.Println("downloading: ", sd.Source.String())

		dest := sd.Destination.Path
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			makePath(dest)
			fmt.Println("Creating: ", dest)
			file, err := os.Create(dest)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			dl := s3.GetObjectInput{
				Bucket: aws.String(sd.Source.Host),
				Key:    aws.String(sd.Source.Path),
			}

			numBytes, err := downloader.Download(file, &dl)
			if err != nil {
				fmt.Println("Failed to download file: ", err.Error())
				continue
			}
			fmt.Printf("%06d bytes downloaded\n", numBytes)
			file.Close()
			nextArgs := SrcDest{
				Source: sd.Destination,
			}
			fmt.Println("Moving")
			next <- nextArgs

		}
	}
}
