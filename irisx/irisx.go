package irisx

import (
	"strings"

	"github.com/kataras/iris/v12"
)

func GetRealIP(ctx iris.Context) string {
	ipStr := ctx.GetHeader("X-Forwarded-For")
	if ipStr != "" {
		ips := strings.Split(ipStr, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	ip := ctx.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	return ctx.RemoteAddr()
}
