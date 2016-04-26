package aws

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/c4labs/c4"
)

// var wg sync.WaitGroup

var scratchPath string
var assetPath string
var metadataPath string

func makePath(file string) {
	path := filepath.Dir(file)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func init() {
	scratchPath = "/tmp/.c4/scratch"
	assetPath = "/tmp/.c4/assets"
	metadataPath = "/tmp/.c4/metadata"

	makePath(scratchPath + "/foo.bar")
	makePath(assetPath + "/foo.bar")
	makePath(metadataPath + "/foo.bar")
}

type fileInfoFunc func(path string, id string, fi os.FileInfo)

type Resource interface {
	Pull(fileInfoFunc)
	Push()
	Scope(string)
	Prefix(string)
}

type SrcDest struct {
	Source      url.URL
	Destination url.URL
}

type AWSResource struct {
	srv      *s3.S3
	sess     *session.Session
	mvChan   chan SrcDest
	idChan   chan SrcDest
	dlChan   chan SrcDest
	objects  []*s3.Object
	uploader *s3manager.Uploader
	bucket   string
	prefix   string
}

type partInfo struct {
	ETag       string
	PartNumber uint
}

type partList struct {
	Parts []partInfo
}

type multipartUploadReport struct {
	Bucket          string
	Key             string
	MultipartUpload partList
}

type partUploadReport struct {
	Bucket     string
	Key        string
	PartNumber uint
}

func AddResource(kind string) (r Resource) {
	switch kind {
	case "aws":
		aws := new(AWSResource)
		aws.sess = session.New()
		aws.srv = s3.New(aws.sess)
		aws.mvChan = make(chan SrcDest)
		aws.idChan = make(chan SrcDest)
		aws.dlChan = make(chan SrcDest)
		r = aws
		go file_mover(aws.mvChan)
		go id_file(aws.idChan, aws.mvChan)
		go file_downloader(aws.dlChan, aws.idChan)
	default:
		aws := new(AWSResource)
		aws.sess = session.New()
		aws.srv = s3.New(aws.sess)
		aws.mvChan = make(chan SrcDest)
		aws.idChan = make(chan SrcDest)
		aws.dlChan = make(chan SrcDest)
		r = aws
		go file_mover(aws.mvChan)
		go id_file(aws.idChan, aws.mvChan)
		go file_downloader(aws.dlChan, aws.idChan)
	}
	return
}

func (r *AWSResource) Scope(s string) {
	r.bucket = s
}

func (r *AWSResource) Prefix(prefix string) {
	r.prefix = prefix
}

func (r *AWSResource) select_files(path string) {
	// fmt.Printf("select_files: %s\n", path)
	// path := filename + string(filepath.Separator) + file.Name()
	f, err := os.Lstat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get status for \"%s\": %s\n", path, err)
		os.Exit(1)
	}
	if f.Mode()&os.ModeSocket == os.ModeSocket {
		return
	}
	if f.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to ReadDir: %v\n", err)
			os.Exit(1)
		}
		for _, file := range files {
			r.select_files(path + "/" + file.Name())
		}
	} else {
		r.upload_file(path)
	}
}

func (r *AWSResource) Push() {
	// Create an uploader with the session and custom options
	r.sess.Handlers.Send.PushFront(func(r *request.Request) {
		fmt.Printf("Request: %s\n", r.Operation.Name)
		if r.Operation.Name == "CompleteMultipartUpload" {
			if mur, ok := r.Params.(multipartUploadReport); ok {
				fmt.Printf("Done: %s\n", mur.Key)
			}
		} else {
			if pur, ok := r.Params.(partUploadReport); ok {
				fmt.Printf("Part %d: %s\n", pur.PartNumber, pur.Key)
			}
		}
		// fmt.Printf("Request: %s/%s, Payload: %s", r.ClientInfo.ServiceName, r.Operation.Name, r.Params)
	})
	r.uploader = s3manager.NewUploader(r.sess, func(u *s3manager.Uploader) {
		// u.PartSize = 16 * 1024 * 1024 // 16MB per part
	})
	r.select_files(assetPath)
}

func file_mover(arg chan SrcDest) {
	for {
		cmd := <-arg
		src := cmd.Source.Path
		dest := cmd.Destination.Path

		fmt.Println("moving file: ", src)
		fmt.Println("to: ", dest)

		makePath(dest)
		err := os.Rename(src, dest)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (r *AWSResource) Pull(fn fileInfoFunc) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(r.bucket), // Required
		Prefix: aws.String(r.prefix),
	}
	pageNum := 0

	err := r.srv.ListObjectsPages(params, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		pageNum++
		for _, obj := range page.Contents {
			path := *obj.Key
			src := url.URL{
				Scheme: "s3",
				Host:   r.bucket,
				Path:   path,
			}
			dest := url.URL{
				Scheme: "file",
				Path:   scratchPath + "/" + path[:5] + path[96:],
			}

			dl := SrcDest{
				Source:      src,
				Destination: dest,
			}
			r.dlChan <- dl
		}
		return pageNum <= 3
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func encode(src io.Reader) *c4.ID {
	e := c4.NewIDEncoder()
	_, err := io.Copy(e, src)
	if err != nil {
		panic(err)
	}
	return e.ID()
}

func id_file(in chan SrcDest, out chan SrcDest) {

	for {
		src := <-in
		f, err := os.Open(src.Source.Path)
		if err != nil {
			fmt.Println(err.Error())
		}
		id := encode(f).String()
		f.Close()
		var destUrl url.URL
		destUrl.Scheme = "file"
		destUrl.Path = assetPath + "/" + id[2:4] + "/" + id[4:6] + "/" + id[6:8] + "/" + id
		dest := SrcDest{
			Source:      src.Source,
			Destination: destUrl,
		}
		out <- dest
	}
}
