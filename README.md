# appserver
app server for go

```
s := appserver.NewServer()
s.Use(appserver.Recover(true))
s.Use(appserver.TokenAuth)

s.ConnectMysqlDb(*dsn, 6, 3)
s.ConnectRedis(*redisServer, 20, 20, 240*time.Second)
s.ConnectTokenRedis(*tokenRedisServer, 20, 20, 240*time.Second)

s.Post("/foo", validate, foo)
s.ListenAndServe(*serverAddress, 3*time.Second)

func validate(c *appserver.Context) {
	c.Next()
}

type form struct {
	Field string `json:"field" valid:"required"`
}
func foo(c *appserver.Context) {
	f := &form{}
	if err := c.LoadJSONAndValidate(f); err != nil {
		c.Error(err)
		return
	}
	c.Success(appserver.util{"field":f.Field})
}
```
