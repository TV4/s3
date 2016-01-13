package s3

import (
	"fmt"
	"sync"
)

type Chunk struct {
	b []byte
	sync.RWMutex
	ID int
}

func (mb *Chunk) Read(p []byte) (int, error) {
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

func (mb *Chunk) WriteAt(p []byte, off int64) (int, error) {
	if int64(len(mb.b)) < off+int64(len(p)) {
		mb.RLock() // FIXME(ivarg): remove?
		nb := make([]byte, off+int64(len(p)))
		if n := copy(nb, mb.b); n != len(mb.b) {
			fmt.Println(off, len(p), len(mb.b), int64(len(mb.b)), len(nb), off+int64(len(p)))
			return n, fmt.Errorf("did not copy correct number of bytes: %d instead of expected %d", n, len(mb.b))
		}
		mb.b = nb
		mb.RUnlock()
	}

	dst := mb.b[off:]
	return copy(dst, p), nil
}
