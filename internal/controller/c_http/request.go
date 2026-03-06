package c_http

import (
	"fmt"
	"net/http"
	"strconv"
)

type requestOptions struct {
	requestPrefix string
}

type RequestOption func(options *requestOptions)

// SetRequestPrefix deprecated
func SetRequestPrefix(requestPrefix string) RequestOption {
	return func(rOptions *requestOptions) {
		rOptions.requestPrefix = requestPrefix
	}
}

type Request struct {
	*http.Request
	requestOptions
}

func NewRequest(r *http.Request, opts ...RequestOption) *Request {
	var rOptions requestOptions
	for _, opt := range opts {
		opt(&rOptions)
	}

	return &Request{
		Request:        r,
		requestOptions: rOptions,
	}
}

func (r *Request) HTTPId() (id int, err error) {
	id, err = strconv.Atoi(r.PathValue("id"))
	if err != nil {
		err = fmt.Errorf("invalid URL")
	}
	return
}
