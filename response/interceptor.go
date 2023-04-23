package response

import (
	"net/http"
)

func NewInterceptor() http.ResponseWriter {
	writer := NewWriter()

	return writer
}

// type Interceptor struct {
// 	// statusCode int
// 	wroteHeader   bool
// 	header        http.Header
// 	body          io.ReadWriter
// 	closeNotifyCh chan bool
// }

// func NewInterceptor() *Interceptor {
// 	return &Interceptor{
// 		header: http.Header{},
// 		body:   &bytes.Buffer{},
// 	}
// }

// func (interceptor *Interceptor) Header() http.Header {
// 	return interceptor.header
// }

// func (interceptor *Interceptor) Write(data []byte) (int, error) {
// 	return interceptor.body.Write(data)
// }

// func (interceptor *Interceptor) WriteHeader(statusCode int) {
// 	// interceptor.statusCode = statusCode
// }
