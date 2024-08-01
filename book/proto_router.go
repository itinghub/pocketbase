package book

import (
	"errors"
	"reflect"

	"github.com/labstack/echo/v5"
)

type RouteItem struct {
	Method string
	Path   string
	req    reflect.Type
	resp   reflect.Type
}

type ProtoRouter struct {
	routes map[string]RouteItem
}

type MatchResult struct {
	Err  error
	Item RouteItem
}

const MethodPathSeprator = " "

func protoRouter(items []RouteItem) ProtoRouter {

	routes := make(map[string]RouteItem)
	for _, item := range items {
		routes[item.Method+MethodPathSeprator+item.Path] = item
	}
	return ProtoRouter{routes: routes}
}

func (r *ProtoRouter) Match(c echo.Context) (result MatchResult, matched bool) {
	key := c.Request().Method + MethodPathSeprator + c.Request().URL.Path
	if item, ok := r.routes[key]; ok {
		if cc, ok := c.(*ProtoContext); ok {
			err := cc.PopulateQueryParam(item.req, item.resp)
			return MatchResult{Err: err, Item: item}, true

		} else {
			return MatchResult{Err: errors.New("CustomContext not match"), Item: item}, true
		}
	} else {
		return MatchResult{}, false
	}
}

// if encodedParams := cc.QueryParam(ProtoQueryParamKey); encodedParams != "" {

// 	decodedParams, err := base64.StdEncoding.DecodeString(encodedParams)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid encoded parameters")
// 	}

// } else {
// 	return errors.New("encodedParams not found")
// }
