package appserver

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
	"github.com/juju/errors"
	"gopkg.in/doug-martin/goqu.v3"
)

// Database 数据库处理对象
type Database interface {
	From(cols ...interface{}) *goqu.Dataset
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Transaction 事务处理对象
type Transaction interface {
	Begin() (*goqu.TxDatabase, error)
}

// DbContext 数据库上下文
type DbContext struct {
	db             Database
	redisPool      *redis.Pool
	tokenRedisPool *redis.Pool
	inTransaction  bool
}

// NewDbContext 新建DbContext
func NewDbContext(db Database, redisPool *redis.Pool, tokenRedisPool *redis.Pool) *DbContext {
	return &DbContext{db: db, redisPool: redisPool, tokenRedisPool: tokenRedisPool}
}

// Transaction 开启事务
func (c *DbContext) Transaction(fn func() error) error {
	if !c.inTransaction {
		tx, err := c.db.(Transaction).Begin()
		if err != nil {
			return errors.Annotate(err, "开启事务")
		}
		c.db = tx
		c.inTransaction = true
		return tx.Wrap(fn)
	}
	return fn()
}

// Db 获取数据库实例或事务实例
func (c *DbContext) Db() Database {
	return c.db
}

// Table 获取数据库表实例
func (c *DbContext) Table(name string) *goqu.Dataset {
	return c.db.From(goqu.I(name))
}

// RedisConn 获取一个redis连接
func (c *DbContext) RedisConn(fn func(redis.Conn) error) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	return fn(conn)
}

// TokenRedisConn 获取一个token专用的redis连接
func (c *DbContext) TokenRedisConn(fn func(redis.Conn) error) error {
	conn := c.tokenRedisPool.Get()
	defer conn.Close()
	return fn(conn)
}
