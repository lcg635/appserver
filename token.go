package appserver

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/garyburd/redigo/redis"
)

const (
	// 默认的token有效其
	defaultTokenExpire = 30 * 86400
	// token的key前缀
	tokenNS = "token:"
)

var tokenEntropy = 32

// TokenAuth token校验中间件
func TokenAuth(c *Context) {
	accountID, err := c.AccountID()
	if err != nil {
		c.Error(err)
		return
	}
	if accountID == 0 {
		c.Error(ErrNotAuthorized)
		return
	}
	c.Next()
}

// AccountID 获取当前的用户id
func (c *Context) AccountID() (int64, error) {
	if c.accountID == 0 {
		token, err := c.Token()
		if err != nil {
			return 0, err
		}
		err = c.DbContext().TokenRedisConn(func(conn redis.Conn) error {
			accountID, err := redis.Int64(conn.Do("GET", tokenNS+token))
			if err != nil && err.Error() != "redigo: nil returned" {
				return err
			}
			c.accountID = accountID
			return nil
		})
		if err != nil {
			return 0, err
		}
	}
	return c.accountID, nil
}

// NewToken 生成一个随机的token
func (c *Context) NewToken(uid int64) (string, error) {
	bytes := make([]byte, tokenEntropy)
	_, err := rand.Read(bytes[:cap(bytes)])
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)

	err = c.DbContext().TokenRedisConn(func(conn redis.Conn) error {
		_, err = conn.Do("SETEX", tokenNS+token, defaultTokenExpire, uid)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	return token, nil
}

// Token 获取当前的token
func (c *Context) Token() (string, error) {
	if c.token == "" {
		authHeader := c.r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Token") {
			return "", ErrInvalidAuthorizationHeader
		}
		_, err := base64.URLEncoding.DecodeString(parts[1])
		if err != nil {
			return "", ErrInvalidTokenEncoding
		}
		c.token = string(parts[1])
	}
	return c.token, nil
}

// RefreshToken 更新token的有效期
func (c *Context) RefreshToken(token string, expire int) error {
	return c.DbContext().TokenRedisConn(func(conn redis.Conn) error {
		_, err := conn.Do("EXPIRE", tokenNS+token, expire)
		if err != nil {
			return err
		}
		return nil
	})
}

// RemoveToken 移除一个token
func (c *Context) RemoveToken(token string) error {
	return c.DbContext().TokenRedisConn(func(conn redis.Conn) error {
		_, err := conn.Do("DEL", tokenNS+token)
		if err != nil {
			return err
		}
		return nil
	})
}
