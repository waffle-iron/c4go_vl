package asset

import (
  "crypto/sha512"
  "hash"
  "io"
  "io/ioutil"
  "math/big"
)

// IDEncoder generates C4 Asset IDs.
type IDEncoder struct {
  err error
  h   hash.Hash
  wr  io.Writer
}

// NewIDEncoder makes a new IDEncoder.
func NewIDEncoder(w_op ...io.Writer) *IDEncoder {
  w := ioutil.Discard

  if len(w_op) > 0 {
    w = w_op[0]
  }

  return &IDEncoder{
    wr: w,
    h:  sha512.New(),
  }
}

// Write writes bytes to the hash that makes up the ID.
func (e *IDEncoder) Write(b []byte) (int, error) {
  _, e.err = e.wr.Write(b)
  return e.h.Write(b)
}

// ID gets the ID for the written bytes.
func (e *IDEncoder) ID() *ID {
  b := new(big.Int)
  b.SetBytes(e.h.Sum(nil))
  id := ID(*b)
  return &id
}
