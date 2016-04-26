package asset

import (
  "bytes"
  "fmt"
  "io"
  "math/big"
  "strconv"
)

// using the flickr character set which removes:
// ['=', '+', '_', '0', 'O', 'I', 'l'] from base64
// to reduce transcription errors, and make friendlier URLs
const (
  charset = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
  base    = 58
)

var (
  lut     [256]byte
  lowbyte = []byte("1")
  prefix  = []byte{'c', '4'}
  idlen   = 90
)

func init() {
  for i := 0; i < len(lut); i++ {
    lut[i] = 0xFF
  }
  for i := 0; i < len(charset); i++ {
    lut[charset[i]] = byte(i)
  }
}

type errBadChar int

func (e errBadChar) Error() string {
  return "non c4 character at position " + strconv.Itoa(int(e))
}

// ID represents a C4 Asset ID.
type ID big.Int

type IDable interface {
  ID() *ID
}

// ID of an ID
func (id *ID) ID() *ID {
  return id
}

func (i *ID) Sum(j *ID) *ID {
  var ids [2]*ID
  ids[0] = i
  ids[1] = j
  l := 0
  r := 1

  if ids[r].Cmp(ids[l]) < 0 {
    r = 0
    l = 1
  }
  e := NewIDEncoder()
  _, err := io.Copy(e, bytes.NewReader(append(ids[l].Bytes(), ids[r].Bytes()...)))
  if err != nil {
    fmt.Print(err)
  }
  return e.ID()
}

// ParseID parses a C4 ID string into an ID.
func ParseID(src string) (*ID, error) {
  return ParseBytesID([]byte(src))
}

// ParseBytesID parses a C4 ID as []byte into an ID.
func ParseBytesID(src []byte) (*ID, error) {
  bigNum := new(big.Int)
  bigBase := big.NewInt(base)
  for i := 2; i < len(src); i++ {
    b := lut[src[i]]
    if b == 0xFF {
      return nil, errBadChar(i)
    }
    bigNum.Mul(bigNum, bigBase)
    bigNum.Add(bigNum, big.NewInt(int64(b)))
  }
  id := ID(*bigNum)
  return &id, nil
}

func (id *ID) String() string {
  return string(id.Bytes())
}

func (x *ID) Cmp(y *ID) int {
  bigX := big.Int(*x)
  bigY := big.Int(*y)
  return bigX.Cmp(&bigY)
}

// Bytes encodes the written bytes to C4 ID format.
func (id *ID) Bytes() []byte {
  var bigNum big.Int
  bigID := big.Int(*id)
  bigNum.Set(&bigID)
  bigBase := big.NewInt(base)
  bigZero := big.NewInt(0)
  bigMod := new(big.Int)
  var encoded []byte
  for bigNum.Cmp(bigZero) > 0 {
    bigNum.DivMod(&bigNum, bigBase, bigMod)
    encoded = append([]byte{charset[bigMod.Int64()]}, encoded...)
  }
  // padding
  diff := idlen - 2 - len(encoded)
  encoded = append(bytes.Repeat(lowbyte, diff), encoded...)
  // c4... prefix
  encoded = append(prefix, encoded...)
  return encoded
}

func (id *ID) Less(idArg *ID) bool {
  {
    return id.Cmp(idArg) < 0
  }
}
