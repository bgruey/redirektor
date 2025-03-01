package qrcode

type WriteCloser struct {
	buf []byte
	n   int
}

// Create new with zero length slice
func NewWriteCloser() *WriteCloser {
	return &WriteCloser{buf: []byte(nil)}
}

// No error checking for writing
func (wc *WriteCloser) Write(p []byte) (n int, err error) {
	wc.buf = append(wc.buf, p...)
	wc.n += len(p)
	return len(p), nil
}

func (wc *WriteCloser) Close() error {
	return nil
}

func (wc *WriteCloser) Bytes() []byte {
	return wc.buf[:wc.n]
}
