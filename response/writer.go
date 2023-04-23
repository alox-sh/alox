package response

import (
	"bytes"
	"io"
	"net/http"
)

type Writer struct {
	wroteHeader   bool
	statusCode    int
	header        http.Header
	body          io.ReadWriter
	closeNotifyCh chan bool
}

func NewWriter() *Writer {
	return &Writer{
		header:        http.Header{},
		body:          &bytes.Buffer{},
		closeNotifyCh: make(chan bool),
	}
}

func (writer *Writer) Header() http.Header {
	if writer.header == nil {
		writer.header = make(http.Header)
	}

	return writer.header
}

func (writer *Writer) Write(data []byte) (int, error) {
	if !writer.wroteHeader {
		writer.WriteHeader(http.StatusOK)
	}

	return writer.body.Write(data)
}

func (writer *Writer) WriteHeader(statusCode int) {
	// writer.statusCode = statusCode

	if writer.wroteHeader {
		return
	}

	if writer.Header().Get("Content-Type") == "" {
		writer.Header().Set("Content-Type", "text/plain; charset=utf8")
	}

	writer.statusCode = statusCode
	// writer.out.StatusCode = statusCode

	// h := make(map[string]string)
	// mvh := make(map[string][]string)

	// for k, v := range writer.Header() {
	// 	if len(v) == 1 {
	// 		h[k] = v[0]
	// 	} else if len(v) > 1 {
	// 		mvh[k] = v
	// 	}
	// }

	// writer.out.Headers = h
	// writer.out.MultiValueHeaders = mvh
	writer.wroteHeader = true
}

func (writer *Writer) CloseNotify() <-chan bool {
	return writer.closeNotifyCh
}
