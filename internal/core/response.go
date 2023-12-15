package core

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RouteResponse struct {
	Type       RouteResponseType
	StatusCode int
	Body       string
	Headers    map[string]string
	Delay      time.Duration
}

func (rr RouteResponse) BuildResponseBody(r *http.Request) string {
	if rr.Type == RESPONSE_TYPE_STATIC {
		return rr.Body
	}
	if rr.Type == RESPONSE_TYPE_DYNAMIC {
		return InterpolateString(rr.Body, r)
	}

	return `{"error":"Not implemented"}`
}

func (rr RouteResponse) BuildResponseStatusCode() int {
	switch rr.Type {
	case RESPONSE_TYPE_STATIC:
		return rr.StatusCode
	case RESPONSE_TYPE_DYNAMIC:
		return rr.StatusCode
	default:
		return http.StatusNotImplemented
	}
}

func (rr RouteResponse) BuildHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for name, value := range rr.Headers {
		fmt.Println(name, value)
		if strings.Contains(value, "$") {
			headers[name] = InterpolateString(value, r)
		} else {
			headers[name] = value
		}
	}
	return headers
}
