package appserver

// 请求方法
const (
	MethodPost = "POST"
	MethodGet  = "GET"
)

// HandlerFunc 请求处理函数
type HandlerFunc func(*Context)

// Route 路由信息
type Route struct {
	Method   string
	Path     string
	handlers []HandlerFunc
}
