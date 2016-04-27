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
	"strings"
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
	recursive_flag 		bool
	version_flag      	bool
	links_flag 			bool
	formatting_string 	string
	depth 				int
	include_meta 		bool
	absolute_flag 		bool
	myBucket          	= []byte("c4bolt_test")
	dbLocation        	string
)

func init() {
	message := "\n" + versionString() + "\n\nUsage: c4bolt [flags] [filepath/url]\n" +
		"  c4bolt generates c4id for specified filepath or url.\n" + "flags:\n"
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, message)
		flag.PrintDefaults()
	}
	flag.BoolVarP(&version_flag, "version", "v", false, "Show version information.")
	flag.BoolVarP(&recursive_flag, "recursive", "R", false, "Recursively identify all files for the given url.")
	flag.BoolVarP(&absolute_flag, "absolute", "a", false, "Output absolute paths, instead of relative paths.")
	flag.BoolVarP(&links_flag, "links", "L", false, "All symbolic links are followed.")
	flag.IntVarP(&depth, "depth", "d", 0, "Only output ids for files and folders 'depth' directories deep.")
	flag.BoolVarP(&include_meta, "metadata", "m", false, "Include filesystem metadata.\n          \"url\" is always included unless data is piped, or only a single file is specified.")
	flag.StringVar(&dbLocation, "db", "mybolt1.db", "Location of your boltdb file")
	flag.StringVarP(&formatting_string, "formatting", "f", "id", "Output formatting options.\n          \"id\": c4id oriented.\n          \"path\": path oriented.")
}

func versionString() string {
	return `c4 version ` + version_number + ` (` + runtime.GOOS + `)`
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("Unable to proceed [%s]", err)
	}
	os.Exit(1)
}

func encode(src io.Reader) *asset.ID {
	e := asset.NewIDEncoder()
	_, err := io.Copy(e, src)
	handleErr(err)
	return e.ID()
}

func fileID(path string) (id *asset.ID) {
	f, err := os.Open(path)
	handleErr(err)
	encode(f)
	f.Close()
	return
}

func nullId() *asset.ID {
	e := asset.NewIDEncoder()
	io.Copy(e, strings.NewReader(``))
	return e.ID()
}
func printID(id *asset.ID) {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf("Id: %s \n\t%s \n", id.String())
	} else {
		fmt.Printf("Id: %s \n\t%s \n", id.String())
	}
}

func newItem(path string) (item map[string]interface{}) {
	item = make(map[string]interface{})
	if item == nil {
		fmt.Fprintf(os.Stderr, "Unable to allocate space for file information for \"%s\".", path)
		os.Exit(1)
	}
	f, err := os.Lstat(path)
	handleErr(err)

	item["folder"] = f.IsDir()
	item["link"] = f.Mode()&os.ModeSymlink == os.ModeSymlink

	return item
}

func generate_c4id(fileurl string, depth int, relative_path string) (id *asset.ID) {
	path, err := filepath.Abs(fileurl)
	handleErr(err)

	item := newItem(path)
	if item["link"] == true && !links_flag {
		newFileUrl, _ := filepath.EvalSymlinks(fileurl)
		item["link"] = newFileUrl
		id = nullId()
	} else {
			if item["link"] == true {
				newFileUrl, err := filepath.EvalSymlinks(fileurl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to follow link %s. %s\n", newFileUrl, err)
				item["link"] = newFileUrl
				id = nullId()
			} else {
				item["link"] = newFileUrl
				var linkId asset.IDSlice
				linkId.Push(generate_c4id(newFileUrl, depth-1, relative_path))
				id = linkId.ID()
			}		
		} 
	}

	item["c4id"] = id.String()
	if depth >=0 || recursive_flag {
		//output(path, item)
		fmt.Printf("Entered output section.........")
	}
	return
}

func generate_c4id_c4id(cid string) (id *asset.ID){
	path, err := filepath.Abs(cid)
	handleErr(err)

	item := newItem(path)
	if item["link"] == true {
		newFileUrl, err := filepath.EvalSymlinks(cid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to follow link %s. %s\n", newFileUrl, err)
			item["link"] = newFileUrl
			id = nullId()
		} else {
			item["link"] = newFileUrl
			var linkId asset.IDSlice
			linkId.Push(generate_c4id_c4id(newFileUrl))
			id = linkId.ID()
		}
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
	fmt.Printf("\nInserted [%s]:[%s] ", id, path)
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
    if version_flag {
		fmt.Println(versionString())
		os.Exit(0)
	}

    if len(fileurl) == 0 {
		flag.Usage()
	} else if len(fileurl) == 1 && !(recursive_flag || include_meta) && depth == 0 {
		path, err := filepath.Abs(fileurl[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to find absolute path for %s. %s\n", fileurl[0], err)
			os.Exit(1)
		}
		id := generate_c4id(path, -1, "")
		printID(id)
	} else {
		for _, file := range fileurl {
			path, err := filepath.Abs(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to find absolute path for %s. %s\n", file, err)
				os.Exit(1)
			}
			if depth < 0 {
				depth = 0
			}

			// Generate c4id & insert to db
			id := generate_c4id(path, depth, "")
			printID(id)
			insert(id.String(), path)

			// Generate c4id for above c4id & insert to db
			id1 := generate_c4id_c4id(id.String())
			printID(id1)
			insert(id1.String(), id.String())
	
			// Read all key:value pairs frm db
			read()		
		}		
	}
	os.Exit(0)
}
