package lambda

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"alox.sh"
)

type APIGatewayHandler func(*APIGateway, context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type APIGateway struct {
	*Lambda
	handler APIGatewayHandler
}

func NewAPIGateway(handler APIGatewayHandler) (apiGateway *APIGateway) {
	apiGateway = &APIGateway{Lambda: NewLambda(apiGateway.lambdaHandler)}

	apiGateway.SetHandler(handler)
	apiGateway.Server.SetHandler(apiGateway.httpHandler)

	return apiGateway
}

func (apiGateway *APIGateway) lambdaHandler(lambda *Lambda, context context.Context, request []byte) (response []byte, err error) {
	var (
		eventRequest  = events.APIGatewayProxyRequest{}
		eventResponse events.APIGatewayProxyResponse
	)

	if err = json.Unmarshal(request, &eventRequest); err != nil {
		return
	}

	if eventResponse, err = apiGateway.handler(apiGateway, context, eventRequest); err != nil {
		return
	}

	return json.Marshal(eventResponse)
}

func (apiGateway *APIGateway) httpHandler(_ alox.Server, responseWriter http.ResponseWriter, request *http.Request) {
	apiGatewayRequest := events.APIGatewayProxyRequest{
		// Resource                        string                        `json:"resource"` // The resource path defined in API Gateway
		// Path                            string                        `json:"path"`     // The url path for the caller
		HTTPMethod:        request.Method,
		Headers:           map[string]string{},
		MultiValueHeaders: map[string][]string{},
		// QueryStringParameters           map[string]string             `json:"queryStringParameters"`
		// MultiValueQueryStringParameters map[string][]string           `json:"multiValueQueryStringParameters"`
		// PathParameters                  map[string]string             `json:"pathParameters"`
		// StageVariables                  map[string]string             `json:"stageVariables"`
		// RequestContext                  APIGatewayProxyRequestContext `json:"requestContext"`
		// Body                            string                        `json:"body"`
		// IsBase64Encoded                 bool                          `json:"isBase64Encoded,omitempty"`
	}

	for key, values := range request.Header {
		switch len(values) {
		case 0:
			apiGatewayRequest.Headers[key] = ""
		case 1:
			apiGatewayRequest.Headers[key] = values[0]
		default:
			apiGatewayRequest.MultiValueHeaders[key] = values
		}
	}

	// TODO: finish

	apiGatewayResponse, err := apiGateway.handler(apiGateway, request.Context(), apiGatewayRequest)
	if err != nil {
		apiGateway.HandleError(responseWriter, request, err)
		return
	}

	header := responseWriter.Header()

	for key, value := range apiGatewayResponse.Headers {
		header.Set(key, value)
	}

	for key, values := range apiGatewayResponse.MultiValueHeaders {
		for _, value := range values {
			header.Add(key, value)
		}
	}

	responseWriter.WriteHeader(apiGatewayResponse.StatusCode)

	if apiGatewayResponse.IsBase64Encoded {
		buffer := make([]byte, base64.StdEncoding.DecodedLen(len(apiGatewayResponse.Body)))

		if _, err = base64.StdEncoding.Decode(buffer, []byte(apiGatewayResponse.Body)); err != nil {
			apiGateway.HandleError(responseWriter, request, err)
			return
		}

		responseWriter.Write(buffer)
		return
	}

	responseWriter.Write([]byte(apiGatewayResponse.Body))
}

func (apiGateway *APIGateway) SetHandler(handler APIGatewayHandler) *APIGateway {
	apiGateway.handler = handler
	return apiGateway
}

func (apiGateway *APIGateway) WriteJSON(responseWriter http.ResponseWriter, json []byte) {
	// alox.WriteJSON(responseWriter, json)

	
}
