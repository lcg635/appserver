package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/doug-martin/goqu.v3"
	"gopkg.in/tylerb/graceful.v1"
)

// H 哈希
type H map[string]interface{}

// Server 服务器
type Server struct {
	db             Database
	redisPool      *redis.Pool
	tokenRedisPool *redis.Pool
	router         *router
	middlewares    []HandlerFunc
	eventBus       *EventBus
}

// NewServer 新建服务器
func NewServer() *Server {
	return &Server{
		router:      newRouter(),
		middlewares: []HandlerFunc{},
		eventBus:    NewEventBus(false),
	}
}

// Use 添加中间件函数
func (s *Server) Use(handler HandlerFunc) *Server {
	s.middlewares = append(s.middlewares, handler)
	return s
}

// ListenAndServe 启动服务器
func (s *Server) ListenAndServe(addr string, timeout time.Duration) {
	graceful.Run(addr, timeout, s)
}

// ServeHTTP http服务
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewContext(s, w, r, nil)
	s.router.appFunc(context)
}

// Post 创建Post请求的路由
func (s *Server) Post(path string, handlers ...HandlerFunc) {
	s.router.add(&Route{
		Method:   MethodPost,
		Path:     path,
		handlers: handlers,
	})
}

// Get 创建Get请求的路由
func (s *Server) Get(path string, handlers ...HandlerFunc) {
	s.router.add(&Route{
		Method:   MethodGet,
		Path:     path,
		handlers: handlers,
	})
}

// EventBus 获取消息总线
func (s *Server) EventBus() *EventBus {
	return s.eventBus
}

// ConnectMysqlDb 连接mysql
func (s *Server) ConnectMysqlDb(
	dsn string, maxOpenConns, maxIdleConns int) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	s.db = goqu.New("mysql", db)
}

// ConnectRedis 连接redis
func (s *Server) ConnectRedis(
	redisServer string,
	maxIdle, maxActive int, idleTimeout time.Duration) {
	s.redisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		MaxActive:   maxActive,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisServer)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// ConnectTokenRedis 连接token用的redis
func (s *Server) ConnectTokenRedis(
	redisServer string,
	maxIdle, maxActive int, idleTimeout time.Duration) {
	s.tokenRedisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		MaxActive:   maxActive,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisServer)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
