package crypto_helper

import (
	"io"
	"net/http"
)

type DecryptoReader struct {
	io.ReadCloser
}

func (d *DecryptoReader) Read(p []byte) (n int, err error) {
	n, err = d.ReadCloser.Read(p)
	for i := 0; i < n; i++ {
		p[i] = p[i] ^ 0x66
	}
	return
}

type EncryptoResponseWriter struct {
	http.ResponseWriter
}

func (e *EncryptoResponseWriter) Write(p []byte) (int, error) {
	ep := make([]byte, len(p), len(p))
	copy(ep, p)
	for i := 0; i < len(p); i++ {
		p[i] = p[i] ^ 0x66
	}
	return e.ResponseWriter.Write(ep)
}
