package lambda

import (
	"context"
	"net/http"

	"alox.sh"
	"alox.sh/server"
)

type LambdaHandler func(*Lambda, context.Context, []byte) ([]byte, error)

type Lambda struct {
	*server.Server
	handler LambdaHandler
}

func NewLambda(handler LambdaHandler) *Lambda {
	lambda := &Lambda{Server: server.NewServer()}

	lambda.SetHandler(handler)
	lambda.Server.SetHandler(func(_ alox.Server, responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.WriteHeader(501)
		responseWriter.Write([]byte("Not Implemented"))
	})

	return lambda
}

func (lambda *Lambda) SetHandler(handler LambdaHandler) *Lambda {
	lambda.handler = handler
	return lambda
}

func (lambda *Lambda) Invoke(context context.Context, request []byte) (response []byte, err error) {
	return lambda.handler(lambda, context, request)
}
