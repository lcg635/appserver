# appserver
app server for go

```
  s := server.NewServer()
	s.Use(server.Recover(true))
	s.Use(server.TokenAuth)

	s.ConnectMysqlDb(*dsn, 6, 3)
	s.ConnectRedis(*redisServer, 20, 20, 240*time.Second)
	s.ConnectTokenRedis(*tokenRedisServer, 20, 20, 240*time.Second)

	s.Post("/foo", foo)
```
