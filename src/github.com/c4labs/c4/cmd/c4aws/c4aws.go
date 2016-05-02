package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

)

// StructureDB -- Graph/fs db
// MetaDataDB -- Information
// DataDB -- Data


type ByteSize int64 const (
  _ = iota  // ignore first value by assigning to blank identifier
  KB ByteSize = 1<<(10*iota)
  MB
  GB
  TB
  PB
  EB
  ZB
  YB
)


type DomainInfo struct {
	ID        c4.ID     `json:"id"`
	Modified  time.Time `json:"modified"`
	CreatedAt time.Time `json:"created_at"`
}

type NodeType int

const (
	FILE    NodeType = 1 + iota
  FOLDER
  LINK
  SOCKET
)

var nodeTypes = [...]string {
	"File",
	"Folder",
	"Link",
	"Socket",
}

func (nodeType NodeType) String() string {
	return nodeTypes[nodeType - 1]
}

type Node struct {
	ID          c4.ID     `json:"id"`
	Up 		      c4.ID		  `json:"parent"` // root id of parent Merkel Tree
	Down        c4.ID     `json:"children"` // root id of children Merkel Tree
	Name			  string	  `json:"name"`
	DomainInfo  c4.ID     `json:"domain_info"` // root id of set of domain specific metadata objects
	Encounter   time.Time `json:"encounter"`
	Type				NodeType
}

type c4fsItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path    	bool      `json:"folder"`
	Link      string    `json:"-"`
	Socket    bool      `json:"-"`
	Bytes     int64     `json:"bytes"`
	Modified  time.Time `json:"modified"`
	CreatedAt time.Time `json:"created_at"`
	// ChildCount   Stack     `json:"-"`
	CurrentCount int `json:"-"`
}

type S3Object struct {
	ETag *string `type:"string"`
	Key *string `min:"1" type:"string"`
	LastModified *time.Time `type:"timestamp" timestampFormat:"iso8601"`
	Owner *Owner `type:"structure"`
	Size *int64 `type:"integer"`
	StorageClass *string `type:"string" enum:"ObjectStorageClass"`
}


func main() {
	awsSession := session.New()
	s3 := s3.New(awsSession)

	params := &s3.ListObjectsInput{
		Bucket: aws.String("test.cccc.io"), // Required
		Prefix: aws.String("0002"),
	}

	pageNum := 0

	err := s3.ListObjectsPages(params, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		pageNum++

    Contents []*Object `type:"list" flattened:"true"`
    	ETag *string `type:"string"`
    	Key *string `min:"1" type:"string"`
    	LastModified *time.Time `type:"timestamp" timestampFormat:"iso8601"`
    	Owner *Owner `type:"structure"`
    	Size *int64 `type:"integer"`
    	StorageClass *string `type:"string" enum:"ObjectStorageClass"`

		CommonPrefixes []*CommonPrefix `type:"list" flattened:"true"`
    Delimiter *string = `type:"string"`
    EncodingType *string = `type:"string" enum:"EncodingType"`
    IsTruncated *bool `type:"boolean"`
    Marker *string `type:"string"`
    MaxKeys *int64 `type:"integer"`
    Name *string `type:"string"`
    NextMarker *string `type:"string"`
    Prefix *string `type:"string"`


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
// func main() {
//   s := c4.NewStack()
//   s.Push(&c4.Node{1})
//   s.Push(&c4.Node{2})
//   s.Push(&c4.Node{3})
//   fmt.Println(s.Pop(), s.Pop(), s.Pop())

//   q := c4.NewQueue(1)
//   q.Push(&c4.Node{4})
//   q.Push(&c4.Node{5})
//   q.Push(&c4.Node{6})
//   fmt.Println(q.Pop(), q.Pop(), q.Pop())
// }


// QQueue<QString> _childIDs;

// func main() {
// 	db, err := bolt.Open(".c4/c4.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer db.Close()
// 	aws := aws.AddResource("aws")
// 	aws.Scope("test.cccc.io")
// 	aws.Prefix("0001/c45ua94E3wZpuPe31tFeA3nMifxQjp47fb9AQ57FykbWxNdRSS9JweayqzWAC3AdLVzrfsqq1FLhAmAmQNyecWsdxv/")
// 	aws.Pull(func(path string, id string, fi os.FileInfo) {
// 		linkPath := ""
// 		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
// 			linkPath, _ = filepath.EvalSymlinks(fi.Name())
// 		}
// 		md := &c4fsItem{
// 			ID:        id,
// 			Name:      fi.Name(),
// 			Folder:    fi.IsDir(),
// 			Link:      linkPath,
// 			Socket:    fi.Mode()&os.ModeSocket == os.ModeSocket,
// 			Bytes:     fi.Size(),
// 			Modified:  fi.ModTime().UTC(),
// 			CreatedAt: time.Now().UTC(),
// 		}
// 		db.Update(func(tx *bolt.Tx) error {
// 			b, err := tx.CreateBucketIfNotExists([]byte("metadata"))
// 			if err != nil {
// 				return err
// 			}
// 			encoded, err := json.Marshal(md)
// 			if err != nil {
// 				return err
// 			}
// 			return b.Put([]byte(path), encoded)
// 		})
// 	})
// 	aws.Push()
// 	// var userRoot map[string]interface{}{
// 	//   "Joshua Kolden <joshua@studiopyxis.com>": map[string]interface{}{
// 	//     "s3": map[string]interface{}{
// 	//       "in.cccc.io": map[string]interface{}
// 	//     }
// 	//   }
// 	// }
// }

// var fsRoot map[string]interface{}

// db.Update(func(tx *bolt.Tx) error {
//   // Assume bucket exists and has keys
//   mdBucket := tx.Bucket([]byte("metadata"))
//   fsBucket := tx.CreateBucketIfNotExists([]byte("fs"))
//   c := b.Cursor()

//   for k, v := c.First(); k != nil; k, v = c.Next() {
//     paths = filepath.SplitList(k)
//     encoded, err := json.Marshal(md)
//     for _, name := range paths {
//       fsBucket.Put([]byte(name), encoded)
//     }
//     fmt.Printf("key=%s, value=%s\n", k, v)
//   }
//   }); err != nil {
//     return err
//   }
//   for k, v := c.First(); k != nil; k, v = c.Next() {
//     fmt.Printf("%s: \n  %s\n", k, v)
//   }
//   return nil
// })
// }
