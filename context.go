package appserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/juju/errors"
)

// Context 每次请求生成这样一个上下文,主要处理事务之类的事情
type Context struct {
	s         *Server
	r         *http.Request
	w         http.ResponseWriter
	token     string
	accountID int64
	index     int
	handlers  []HandlerFunc
	dbContext *DbContext
	nowUnix   int64
	Env       H
}

// NewContext 新建上下文环境
func NewContext(s *Server, w http.ResponseWriter, r *http.Request, handlers []HandlerFunc) *Context {
	return &Context{
		s:        s,
		r:        r,
		w:        w,
		index:    -1,
		handlers: handlers,
		nowUnix:  time.Now().Unix(),
		Env:      make(H),
	}
}

// Request 返回http请求对象
func (c *Context) Request() *http.Request {
	return c.r
}

// ApplyEvent 发出事件
func (c *Context) ApplyEvent(e Event) {
	c.s.EventBus().ApplyEvent(c, e)
}

// Error 发送错误信息
func (c *Context) Error(e error) {
	// logrus.Error(e)
	if err := c.JSON(H{"error": e.Error()}); err != nil {
		logrus.Error(err)
	}
}

// Success 发送成功信息
func (c *Context) Success(v interface{}) {
	if err := c.JSON(H{"error": "", "data": v}); err != nil {
		logrus.Error(err)
	}
}

// JSON 发送JSON响应
func (c *Context) JSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.Header("Content-Type", "application/json")
	if _, err = c.w.Write(b); err != nil {
		return err
	}

	return nil
}

// Header 设置响应头
func (c *Context) Header(key, value string) {
	c.w.Header().Set(key, value)
}

// LoadJSONAndValidate 从请求中加载json并且验证数据
func (c *Context) LoadJSONAndValidate(v interface{}) error {
	if err := c.LoadJSON(v); err != nil {
		return err
	}
	_, err := govalidator.ValidateStruct(v)
	if err != nil {
		return errors.Annotate(err, "请求校验失败")
	}
	return nil
}

// LoadJSON 从请求中加载json并且验证数据
func (c *Context) LoadJSON(v interface{}) error {
	content, err := ioutil.ReadAll(c.r.Body)
	c.r.Body.Close()
	if err != nil {
		return errors.Annotate(err, "解析请求中的json数据失败")
	}
	if len(content) == 0 {
		return ErrJSONPayloadEmpty
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return err
	}
	return nil
}

// NowUnix 上下文的创建时间时间戳
func (c *Context) NowUnix() int64 {
	return c.nowUnix
}

// Next 运行下一个处理函数
func (c *Context) Next() {
	c.index = c.index + 1
	if c.handlers == nil || c.handlers[c.index] == nil {
		return
	}
	c.handlers[c.index](c)
}

// DbContext 获取当前的DbContext
func (c *Context) DbContext() *DbContext {
	if c.dbContext == nil {
		c.dbContext = NewDbContext(c.s.db, c.s.redisPool, c.s.tokenRedisPool)
	}
	return c.dbContext
}
