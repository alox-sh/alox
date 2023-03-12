package middleware

// func Options(server alox.Server, responseWriter http.ResponseWriter, request *http.Request) {
// 	responseWriter.Header().Set("Content-Length", "0")

// 	if request.ContentLength != 0 {
// 		io.Copy(io.Discard, http.MaxBytesReader(responseWriter, request.Body, 4<<10))
// 	}
// }
