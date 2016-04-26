package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/boltdb/bolt"
	"github.com/c4labs/c4/asset"
	flag "github.com/ogier/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

/*
	1) This file creates a db (using boltDB)
	2) Generates c4id for specified file-path-url
	3) Inserts the [c4id:file-path-url] pair to db
	4) Again generates c4id for the above c4id & inserts this [key:value] pair to db
	5) Reads all key:value pairs from DB & prints at console
*/

const version_number = "1.0"

var (
	version_flag      bool
	formatting_string string
	myBucket          = []byte("perftest")
	dbLocation        string
)

func init() {
	message := versionString() + "\n\nUsage: c4id_bolt [flags] [file]\n\n" +
		"  c4id_bolt generates c4id for specified filepath or url.\n" +
		"flags:\n"
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, message)
		flag.PrintDefaults()
	}
	flag.BoolVarP(&version_flag, "version", "v", false, "Show version information.")
	flag.StringVar(&dbLocation, "db", "mybolt.db", "Location of your boltdb file")
	flag.StringVarP(&formatting_string, "formatting", "f", "id", "Output formatting options.\n          \"id\": c4id oriented.\n          \"path\": path oriented.")
}

func versionString() string {
	return `c4 version ` + version_number + ` (` + runtime.GOOS + `)`
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("Unable to proceed [%s]", err)
	}
}

func encode(src io.Reader) *asset.ID {
	e := asset.NewIDEncoder()
	_, err := io.Copy(e, src)
	if err != nil {
		panic(err)
	}
	return e.ID()
}

func fileID(path string) (id *asset.ID) {
	f, err := os.Open(path)
	handleErr(err)
	encode(f)
	f.Close()
	return
}

func printID(path string, id *asset.ID) {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf("Fileurl: %s \n\t%s \n", path, id.String())
	} else {
		fmt.Printf("FileUrl: %s \n\t%s \n", path, id.String())
	}
}

func newItem(path string) (item map[string]interface{}) {
	item = make(map[string]interface{})

	f, err := os.Lstat(path)
	handleErr(err)

	item["folder"] = f.IsDir()
	item["link"] = f.Mode()&os.ModeSymlink == os.ModeSymlink

	return item
}

func generate_c4id(fileurl string) (id *asset.ID) {
	path, err := filepath.Abs(fileurl)
	handleErr(err)

	item := newItem(path)
	if item["link"] == true {
		newFileUrl, err := filepath.EvalSymlinks(fileurl)
		handleErr(err)
		item["link"] = newFileUrl
		var linkId asset.IDSlice
		linkId.Push(generate_c4id(newFileUrl))
		id = linkId.ID()
	} else {
		id = fileID(path)
	}

	item["c4id"] = id.String()
	return
}

func insert(id string, path string) {
	db, err := bolt.Open(dbLocation, 0644, nil)
	handleErr(err)
	defer db.Close()

	// store to data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(myBucket)
		handleErr(err)

		err = bucket.Put([]byte(id), []byte(path))
		handleErr(err)
		return err
	})
	// fmt.Printf("\nInserted [%s]:[%s] ", key, value)
	return
}

func read() {
	db, err := bolt.Open(dbLocation, 0644, nil)
	handleErr(err)
	defer db.Close()
	fmt.Printf("\n\nReading DB")
	db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(myBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("\n\t[%s]:[%s] ", k, v)
		}
		return nil
	})
	return
}

func main() {
	flag.Parse()
	fileurl := flag.Args()
	fmt.Printf("Starting with dbpath [%s]", dbLocation)

	path, err := filepath.Abs(fileurl[0])
	handleErr(err)

	// Generate c4id & insert to db
	id := generate_c4id(path)
	printID(path, id)
	insert(id.String(), path)

	// Generate c4id for above c4id & insert to db
	id1 := generate_c4id(id.String())
	printID(id1.String(), id)
	insert(id1.String(), id.String())

	// Read all key:value pairs frm db
	read()

	os.Exit(0)
}
