package audiotools

import (
	"bytes"
)

// Creating custom struct to fit into WriteSeeker interface
// to pass it into the wav/mp3 encoder later
type FileBuffer struct {
	Buffer bytes.Buffer
}

func NewFileBuffer() *FileBuffer {
	return &FileBuffer{}
}

func (fb *FileBuffer) Bytes() []byte {
	return fb.Buffer.Bytes()
}

func (fb *FileBuffer) Write(p []byte) (int, error) {
	return fb.Buffer.Write(p)
}

// bypass all seek because bytes.Buffer doesn't support Seek
// NOTE good enough for wav and mp3 encoder
func (fb *FileBuffer) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
