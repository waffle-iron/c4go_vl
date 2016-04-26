package main

import (
  "bufio"
  "fmt"
  "github.com/c4labs/c4/asset"
  flag "github.com/ogier/pflag"
  "golang.org/x/crypto/ssh/terminal"
  "io"
  "io/ioutil"
  "os"
  "path/filepath"
  "runtime"
  "strings"
)

const version_number = "0.5"

func versionString() string {
  return `c4 version ` + version_number + ` (` + runtime.GOOS + `)`
}

var recursive_flag bool
var version_flag bool
var arg_links bool
var links_flag bool
var no_links bool
var summery bool
var depth int
var include_meta bool
var absolute_flag bool
var formatting_string string

func init() {
  message := versionString() + "\n\nUsage: c4 [flags] [file]\n\n" +
    "  c4 generates c4ids for all files and folders spacified.\n" +
    "  If no file is given c4 will read piped data.\n" +
    "  Output is in YAML format.\n\n" +
    "flags:\n"
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, message)
    flag.PrintDefaults()
  }
  flag.BoolVarP(&version_flag, "version", "v", false, "Show version information.")
  flag.BoolVarP(&recursive_flag, "recursive", "R", false, "Recursively identify all files for the given url.")
  flag.BoolVarP(&absolute_flag, "absolute", "a", false, "Output absolute paths, instead of relative paths.")
  // flag.BoolVarP(&arg_links, "arg_links", "H", false, "If the -R option is specified, symbolic links on the command line are followed.\n          (Symbolic links encountered in the tree traversal are not followed by default.)")
  flag.BoolVarP(&links_flag, "links", "L", false, "All symbolic links are followed.")
  // flag.BoolVarP(&no_links, "no_links", "P", true, "If the -R option is specified, no symbolic links are followed.  This is the default.")
  flag.IntVarP(&depth, "depth", "d", 0, "Only output ids for files and folders 'depth' directories deep.")
  flag.BoolVarP(&include_meta, "metadata", "m", false, "Include filesystem metadata.\n          \"url\" is always included unless data is piped, or only a single file is specified.")
  flag.StringVarP(&formatting_string, "formatting", "f", "id", "Output formatting options.\n          \"id\": c4id oriented.\n          \"path\": path oriented.")
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

func printID(id *asset.ID) {
  if terminal.IsTerminal(int(os.Stdout.Fd())) {
    fmt.Printf("%s\n", id.String())
  } else {
    fmt.Printf("%s", id.String())
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

func newItem(path string) (item map[string]interface{}) {
  item = make(map[string]interface{})
  if item == nil {
    fmt.Fprintf(os.Stderr, "Unable to allocate space for file information for \"%s\".", path)
    os.Exit(1)
  }
  f, err := os.Lstat(path)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to get status for \"%s\": %s\n", path, err)
    os.Exit(1)
  }

  item["folder"] = f.IsDir()
  item["link"] = f.Mode()&os.ModeSymlink == os.ModeSymlink
  item["socket"] = f.Mode()&os.ModeSocket == os.ModeSocket
  item["bytes"] = f.Size()
  item["modified"] = f.ModTime().UTC()
  item["currentTime"] = time.Time.Now().UTC()

  return item
}

func walkFilesystem(depth int, filename string, relative_path string) (id *asset.ID) {
  path, err := filepath.Abs(filename)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to find absolute path for %s. %s\n", filename, err)
    os.Exit(1)
  }

  item := newItem(path)
  if item["socket"] == true {
    id = nullId()
  } else if item["link"] == true && !links_flag {
    newFilepath, _ := filepath.EvalSymlinks(filename)
    item["link"] = newFilepath
    id = nullId()
  } else if item["link"] == true {
    newFilepath, err := filepath.EvalSymlinks(filename)
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
  } else {
    if item["folder"] == true {
      files, err := ioutil.ReadDir(path)
      if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to ReadDir: %v\n", err)
        os.Exit(1)
      }
      var childIDs asset.IDSlice
      for _, file := range files {
        path := filename + string(filepath.Separator) + file.Name()
        childIDs.Push(walkFilesystem(depth-1, path, relative_path))
      }
      id = childIDs.ID()
    } else {
      id = fileID(path)
    }
  }
  item["c4id"] = id.String()
  if depth >= 0 || recursive_flag {
    output(path, item)
  }
  return
}

func main() {
  flag.Parse()
  file_list := flag.Args()
  if version_flag {
    fmt.Println(versionString())
    os.Exit(0)
  }

  if len(file_list) == 0 {
    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) == 0 {
      reader := bufio.NewReader(os.Stdin)
      printID(encode(reader))
    } else {
      flag.Usage()
    }
  } else if len(file_list) == 1 && !(recursive_flag || include_meta) && depth == 0 {
    path, err := filepath.Abs(file_list[0])
    if err != nil {
      fmt.Fprintf(os.Stderr, "Unable to find absolute path for %s. %s\n", file_list[0], err)
      os.Exit(1)
    }
    id := walkFilesystem(-1, path, "")
    printID(id)
  } else {
    for _, file := range file_list {
      path, err := filepath.Abs(file)
      if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to find absolute path for %s. %s\n", file, err)
        os.Exit(1)
      }
      if depth < 0 {
        depth = 0
      }
      walkFilesystem(depth, path, "")
    }
  }
  os.Exit(0)
}
