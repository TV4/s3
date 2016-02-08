package s3

import (
	"fmt"
	"io"
)

// Object is a writable structure representing a binary blob residing on an S3
// bucket. It is written as part of the process of downloading objects from
// S3.
type Object struct {
	b  []byte
	ID int
}

// Read reads the next len(p) bytes from the object or until the object is
// fully read. The return value is the number of bytes read. If the object
// has no data, err is io.EOF.
func (mb *Object) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	if len(mb.b) == 0 {
		return 0, io.EOF
	}

	n := len(mb.b)
	if len(p) < n {
		n = len(p)
	}

	nn := copy(p, mb.b[:n])
	if nn < n {
		n = nn
	}
	mb.b = mb.b[n:]

	return n, nil
}

// WriteAt writes the buffer p to the object buffer, starting at position off.
// The return value is the number of bytes written.
func (mb *Object) WriteAt(p []byte, off int64) (int, error) {
	if int64(len(mb.b)) < off+int64(len(p)) {
		nb := make([]byte, off+int64(len(p)))
		if n := copy(nb, mb.b); n != len(mb.b) {
			return n, fmt.Errorf("copied %d bytes instead of the expected %d", n, len(mb.b))
		}
		mb.b = nb
	}

	dst := mb.b[off:]
	return copy(dst, p), nil
}
