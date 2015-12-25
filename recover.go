package appserver

import (
	"fmt"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
)

// Recover 错误恢复的中间件
func Recover(enableResponseStackTrace bool) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if reco := recover(); reco != nil {
				trace := debug.Stack()
				// log the trace
				message := fmt.Sprintf("%s\n%s", reco, trace)
				logrus.Error(message)

				if enableResponseStackTrace {
					c.Error(fmt.Errorf("%s", message))
				} else {
					c.Error(ErrInternalServer)
				}
			}
		}()
		c.Next()
	}
}
