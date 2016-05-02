package main

import (
  "bufio"
  "fmt"
  "time"
  "github.com/boltdb/bolt"
  "github.com/c4labs/c4/asset"
  flag "github.com/ogier/pflag"
  "golang.org/x/crypto/ssh/terminal"
  "io"
  "io/ioutil"
  "net/url"
  "os"
  "path/filepath"
  "runtime"
  "strings"
)

const version_number = "1.0"

func versionString() string {
  return `c4 version ` + version_number + ` (` + runtime.GOOS + `)`
}

var (
	recursive_flag bool
	version_flag bool
	arg_links bool
	links_flag bool
	no_links bool
	summary bool
	depth int
	include_meta bool
	absolute_flag bool
	formatting_string string
	myBucket = []byte("c4bolt_test")
	dbLocation string
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
	flag.StringVar(&dbLocation, "db", "mybolt.db", "Location of your boltdb file")
	flag.StringVarP(&formatting_string, "formatting", "f", "id", "Output formatting options.\n          \"id\": c4id oriented.\n          \"path\": path oriented.")
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

func encodestr(src string) *asset.ID {
  	e := asset.NewIDEncoder()
  	io.Copy(e, strings.NewReader(src))
  	return e.ID()
}

func fileID(path string) (id *asset.ID) {
  f, err := os.Open(path)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to identify %s. %v\n", path, err)
    os.Exit(1)
  }
  id = encode(f)
  f.Close()
  return
}

func nullId() *asset.ID {
  e := asset.NewIDEncoder()
  io.Copy(e, strings.NewReader(``))
  return e.ID()
}

func printID(path string, id *asset.ID) {
  if terminal.IsTerminal(int(os.Stdout.Fd())) {
    fmt.Printf("\n%s : %s", id.String(), path)
  } else {
    fmt.Printf("\n%s : %s", id.String(), path)
  }
}

func output(path string, item map[string]interface{}) {
  rootPath, _ := filepath.Abs(".")
  baseName := filepath.Base(path)

  newPath := path
  if !absolute_flag {
    newPath, _ = filepath.Rel(rootPath, path)
  }
  if include_meta {
    if formatting_string == "path" {
      fmt.Printf("\"%s\":\n", newPath)
      fmt.Printf("  c4id: %s\n", item["c4id"])
    } else {
      fmt.Printf("%s:\n", item["c4id"])
      fmt.Printf("  path: \"%s\"\n", newPath)
    }
    fmt.Printf("  name:  \"%s\"\n", baseName)
    if item["folder"] == false {
      fmt.Printf("  folder:  false\n")
    } else {
      fmt.Printf("  folder:  true\n")
    }
    if item["link"] == false {
      fmt.Printf("  link:  false\n")
    } else {
      linkPath := item["link"].(string)
      if !absolute_flag {
        linkPath, _ = filepath.Rel(rootPath, linkPath)
      }
      fmt.Printf("  link:  \"%s\"\n", linkPath)
    }
    fmt.Printf("  bytes:  %d\n", item["bytes"])
  } else {
    if formatting_string == "path" {
      fmt.Printf("%s:  %s\n", newPath, item["c4id"])
    } else {
      fmt.Printf("%s:  %s\n", item["c4id"], newPath)
    }
  }
}

func isValidUrl(urlStr string)(flag bool){
	_, err := url.Parse(urlStr)
	var validUrl bool
	if err != nil {
	    fmt.Fprintf(os.Stderr, "Unable to get status for \"%s\": %s\n", urlStr, err)
        validUrl = false
    } else {
		validUrl = true
    }
    return validUrl
}

func newItem(path string) (item map[string]interface{}) {
	item = make(map[string]interface{})
	if item == nil {
		fmt.Fprintf(os.Stderr, "Unable to allocate space for file information for \"%s\".", path)
	    os.Exit(1)
	}
	f, err := os.Lstat(path)
	if err != nil {
		if !isValidUrl(path) {
    		fmt.Fprintf(os.Stderr, "Unable to get status for \"%s\": %s\n", path, err)
    		os.Exit(1)
    	} else {
			fmt.Printf("\nReturning = %s \n", path)
			return
    	}
	}
  
  	item["folder"] = f.IsDir()
  	item["link"] = f.Mode()&os.ModeSymlink == os.ModeSymlink
  	item["socket"] = f.Mode()&os.ModeSocket == os.ModeSocket
  	item["bytes"] = f.Size()
  	item["modified"] = f.ModTime().UTC()
  	item["currentTime"] = time.Now().UTC()

  	return item
}

func insert(id string, path string) {
	db, err := bolt.Open(dbLocation, 0600, nil)
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
	handleErr(err)
	fmt.Printf("\nInserting to Db := [%s]:[%s] ", id, path)
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

func walkFilesystem(depth int, fileurl string, relative_path string) (id *asset.ID) {
  	path, err := filepath.Abs(fileurl)
  	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find absolute path for %s. %s\n", fileurl, err)
    	os.Exit(1)
  	}

  	item := newItem(path)
  	if item["socket"] == true {
    	id = nullId()
  	} else if item["link"] == true && !links_flag {
    	newFilepath, _ := filepath.EvalSymlinks(fileurl)
    	item["link"] = newFilepath
    	id = nullId()
  	} else if item["link"] == true {
    	newFilepath, err := filepath.EvalSymlinks(fileurl)
    	if err != nil {
      		fmt.Fprintf(os.Stderr, "Unable to follow link %s. %s\n", newFilepath, err)
      		item["link"] = newFilepath
      		id = nullId()
    	} else {
      		item["link"] = newFilepath
      		var linkId asset.IDSlice
      		linkId.Push(walkFilesystem(depth-1, newFilepath, relative_path))
      		id = linkId.ID()
    	}
  	} else if isValidUrl(path) {
  		id = encodestr(path)
  	} else {
    	if item["folder"] == true {
      		files, err := ioutil.ReadDir(path)
      		if err != nil {
        		fmt.Fprintf(os.Stderr, "Unable to ReadDir: %v\n", err)
        		os.Exit(1)
      		}
      		var childIDs asset.IDSlice
      		for _, file := range files {
        		path := fileurl + string(filepath.Separator) + file.Name()
        		childIDs.Push(walkFilesystem(depth-1, path, relative_path))
      		}
      	id = childIDs.ID()
    	} else {
      		id = fileID(path)
    	}
  	}
  	item["c4id"] = id.String()
  	/*
  	if depth >= 0 || recursive_flag {
    	output(path, item)
  	}*/
  	return
}

func main() {
  flag.Parse()
  file_list := flag.Args()
  if version_flag {
    fmt.Println(versionString())
    os.Exit(0)
  }
  fmt.Printf("Environment = %s/%s \n", os.Getenv("GOOS"), os.Getenv("GOARCH"))

  if len(file_list) == 0 {
    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) == 0 {
      reader := bufio.NewReader(os.Stdin)
      printID("",encode(reader))
    } else {
      flag.Usage()
    }
  } else if len(file_list) == 1 && !(recursive_flag || include_meta) && depth == 0 {
  	if isValidUrl(file_list[0]) {
	    fmt.Printf("Url = %s", file_list[0])  	
	    id := encodestr(file_list[0])

	    fmt.Printf("\nOpening boltdb = %s\\%s", os.Getenv("GOPATH"), dbLocation)
    	insert(id.String(), file_list[0])

		// Generate c4id for above c4id
		fmt.Printf("\n\nGenerating C4Id for above C4Id,")
		id1 := encodestr(id.String())
		printID(id.String(), id1)
	
		insert(id1.String(), id.String())
  	} else {
	    path, err := filepath.Abs(file_list[0])
	    if err != nil {
	      fmt.Fprintf(os.Stderr, "Unable to find absolute path for %s. %s\n", file_list[0], err)
	      os.Exit(1)
	    }
	    id := walkFilesystem(-1, path, "")
	    fmt.Printf("\nOpening boltdb = %s\\%s", os.Getenv("GOPATH"), dbLocation)
    	insert(id.String(), path)

		// Generate c4id for above c4id
		fmt.Printf("\n\nGenerating C4Id for above C4Id,")
		id1 := encodestr(id.String())
		printID(id.String(), id1)
		
		insert(id1.String(), id.String())
	}
    // Read all key:value pairs frm db
	read()
  } else {
  	fmt.Printf("\n\nOpening boltdb = %s\\%s \n", os.Getenv("GOPATH"), dbLocation)
    for _, file := range file_list {
      path, err := filepath.Abs(file)
      if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to find absolute path for %s. %s\n", file, err)
        os.Exit(1)
      }
      if depth < 0 {
        depth = 0
      }
      id := walkFilesystem(depth, path, "")	  
	  insert(id.String(), path)

	  // Generate c4id for above c4id
	  fmt.Printf("\nGenerating C4Id for above C4Id,")
	  id1 := encodestr(id.String())
	  printID(id.String(), id1)
 	
	  insert(id1.String(), id.String())
    }
    // Read all key:value pairs frm boltDb
	read()
  }
  os.Exit(0)
}
