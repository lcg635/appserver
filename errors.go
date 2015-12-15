package server

import "errors"

// 定义一些错误信息
var (
	ErrNotAuthorized              = errors.New("Not Authorized")
	ErrRouteNotFound              = errors.New("Resource Not Found")
	ErrJSONPayloadEmpty           = errors.New("请求中的json数据为空")
	ErrInvalidAuthorizationHeader = errors.New("Invalid Authorization header")
	ErrInvalidTokenEncoding       = errors.New("Token encoding not valid")
	ErrInternalServer             = errors.New("Internal Server Error")
)
