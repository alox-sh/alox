package filter

import (
	"net/http"
	"testing"
)

func TestMethods(test *testing.T) {
	methodsFilter := Methods("put", "bar")

	if methodsFilter(&http.Request{Method: "get"}) {
		test.Fail()
	}

	if methodsFilter(&http.Request{Method: "foo"}) {
		test.Fail()
	}

	if !methodsFilter(&http.Request{Method: "put"}) {
		test.Fail()
	}

	if !methodsFilter(&http.Request{Method: "bar"}) {
		test.Fail()
	}
}
