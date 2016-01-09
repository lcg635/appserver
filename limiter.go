package appserver

import "strings"

// IPLimiter 限制ip访问
func IPLimiter(ips ...string) HandlerFunc {
	limitedIPs := map[string]bool{}
	for _, ip := range ips {
		limitedIPs[ip] = true
	}
	return func(c *Context) {
		ip := c.r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = strings.Split(c.r.RemoteAddr, ":")[0]
		}
		if _, ok := limitedIPs[ip]; !ok {
			c.Error(ErrIPLimited)
			return
		}
		c.Next()
	}
}
