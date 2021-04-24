package util

import (
	"bytes"
	"io"
)

// ReaderToByte reads an io.Reader into a byte slice and returns it.
func ReaderToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
