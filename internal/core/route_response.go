package core

import (
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

func NewRouteResponse(t RouteResponseType, statusCode int, body string, headers map[string]string, delay time.Duration) *RouteResponse {
	return &RouteResponse{
		Type:       t,
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
		Delay:      delay,
	}
}

func (rr RouteResponse) BuildResponseBody(r *http.Request) (*string, error) {
	switch rr.Type {
	case RESPONSE_TYPE_STATIC:
		return &rr.Body, nil
	case RESPONSE_TYPE_DYNAMIC:
		return ProcessString(r, rr.Body)
	default:
		return nil, ErrResponseNotImplemented
	}
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

func (rr RouteResponse) BuildHeaders(r *http.Request) (map[string]string, error) {
	headers := make(map[string]string)
	for name, value := range rr.Headers {
		if strings.Contains(value, "{{") {
			val, err := ProcessString(r, value)
			if err != nil {
				return nil, err
			}
			headers[name] = *val
		} else {
			headers[name] = value
		}
	}
	return headers, nil
}
