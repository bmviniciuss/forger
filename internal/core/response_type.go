package core

type RouteResponseType string

const (
	RESPONSE_TYPE_STATIC  RouteResponseType = "STATIC"
	RESPONSE_TYPE_DYNAMIC RouteResponseType = "DYNAMIC"
)

func NewRouteResponseType(t string) (RouteResponseType, error) {
	switch t {
	case RESPONSE_TYPE_STATIC.String():
		return RESPONSE_TYPE_STATIC, nil
	case RESPONSE_TYPE_DYNAMIC.String():
		return RESPONSE_TYPE_DYNAMIC, nil
	default:
		return "", ErrInvalidRouteResponseType
	}
}

func (rr RouteResponseType) String() string {
	return string(rr)
}
