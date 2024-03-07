package core

type RouteDefinition struct {
	Path     string
	Method   string
	Response RouteResponse
}

func NewRouteDefinition(path, method string, response RouteResponse) *RouteDefinition {
	return &RouteDefinition{path, method, response}
}
