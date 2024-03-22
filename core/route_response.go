package core

import (
	"io"
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

type Result struct {
	StatusCode int
	Body       *string
	Headers    map[string]string
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

func (rr RouteResponse) BuildResponse(r *http.Request) (Result, error) {
	reqBody, err := readBody(r)
	if err != nil {
		return Result{}, err
	}

	body, err := rr.buildResponseBody(r, &reqBody)
	if err != nil {
		return Result{}, err
	}
	headers, err := rr.buildHeaders(r, &reqBody)
	if err != nil {
		return Result{}, err
	}
	return Result{
		StatusCode: rr.buildResponseStatusCode(),
		Body:       body,
		Headers:    headers,
	}, nil
}

func readBody(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}
	rawReqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	return string(rawReqBody), nil
}

func (rr RouteResponse) buildResponseBody(r *http.Request, reqBody *string) (*string, error) {
	switch rr.Type {
	case RESPONSE_TYPE_STATIC:
		return &rr.Body, nil
	case RESPONSE_TYPE_DYNAMIC:
		return processString(r, rr.Body, reqBody)
	default:
		return nil, ErrResponseNotImplemented
	}
}

func (rr RouteResponse) buildResponseStatusCode() int {
	switch rr.Type {
	case RESPONSE_TYPE_STATIC:
		return rr.StatusCode
	case RESPONSE_TYPE_DYNAMIC:
		return rr.StatusCode
	default:
		return http.StatusNotImplemented
	}
}

func (rr RouteResponse) buildHeaders(r *http.Request, reqBody *string) (map[string]string, error) {
	headers := make(map[string]string)
	for name, value := range rr.Headers {
		if strings.Contains(value, "{{") {
			val, err := processString(r, value, reqBody)
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
