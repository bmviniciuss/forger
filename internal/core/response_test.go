package core

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ptrValue[T any](val T) *T {
	return &val
}

func Test_Response_BuildResponseBody(t *testing.T) {
	t.Run("should return static body", func(t *testing.T) {
		rr := RouteResponse{
			Type: RESPONSE_TYPE_STATIC,
			Body: "static body",
		}
		body, err := rr.BuildResponseBody(nil)
		assert.Nil(t, err)
		assert.Equal(t, &rr.Body, body)
	})

	t.Run("should return ErrResponseNotImplemented if RouteResponseType is invalid", func(t *testing.T) {
		rr := RouteResponse{
			Type: RouteResponseType("any"),
			Body: "static body",
		}
		body, err := rr.BuildResponseBody(nil)
		assert.Nil(t, body)
		assert.Equal(t, ErrResponseNotImplemented, err)
	})

	t.Run("should generate uuid in dynamic response", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		rr := RouteResponse{
			Type: RESPONSE_TYPE_DYNAMIC,
			Body: `{{ uuid }}`,
		}
		body, err := rr.BuildResponseBody(r)
		assert.Nil(t, err)
		assert.NotNil(t, body)
		_, err = uuid.Parse(*body)
		assert.Nil(t, err)
	})

	t.Run("should generate uuid of type ulid in dynamic response", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		rr := RouteResponse{
			Type: RESPONSE_TYPE_DYNAMIC,
			Body: `{{ uuid "ulid" }}`,
		}
		body, err := rr.BuildResponseBody(r)
		assert.Nil(t, err)
		assert.NotNil(t, body)
		_, err = uuid.Parse(*body)
		assert.Nil(t, err)
	})

	t.Run("should generate uuid v4 in dynamic response", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		rr := RouteResponse{
			Type: RESPONSE_TYPE_DYNAMIC,
			Body: `{{ uuid "v4" }}`,
		}
		body, err := rr.BuildResponseBody(r)
		assert.Nil(t, err)
		assert.NotNil(t, body)
		_, err = uuid.Parse(*body)
		assert.Nil(t, err)
	})
}
