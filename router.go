package appserver

import "net/url"

// Router 路由控制器
type router struct {
	getRoutes  map[string][]HandlerFunc
	postRoutes map[string][]HandlerFunc
}

func newRouter() *router {
	return &router{
		getRoutes:  make(map[string][]HandlerFunc),
		postRoutes: make(map[string][]HandlerFunc),
	}
}

// AppFunc Handle the REST routing and run the user code.
func (rt *router) appFunc(context *Context) {
	// find the route
	handlers := rt.findHandlersFromURL(context.r.Method, context.r.URL)
	if handlers == nil || len(handlers) == 0 {
		context.Error(ErrRouteNotFound)
		return
	}
	context.handlers = append(context.s.middlewares, handlers...)
	context.Next()
}

// 添加一个路由规则
func (rt *router) add(route *Route) {
	if route.Method == MethodPost {
		rt.postRoutes[route.Path] = route.handlers
	} else if route.Method == MethodGet {
		rt.getRoutes[route.Path] = route.handlers
	}
}

// Return the first matching Route and the corresponding parameters for a given URL object.
func (rt *router) findHandlersFromURL(method string, urlObj *url.URL) []HandlerFunc {
	var handlers []HandlerFunc
	if method == MethodPost {
		handlers = rt.postRoutes[urlObj.Path]
	} else if method == MethodGet {
		handlers = rt.getRoutes[urlObj.Path]
	}
	return handlers
}
