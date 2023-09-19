package irisx

import (
	"errors"
	"strconv"
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

func GetUid(ctx iris.Context) (int64, error) {
	u := ctx.User()
	if u == nil {
		return 0, errors.New("invalid access_token")
	}
	id, err := ctx.User().GetID()
	if err != nil {
		return 0, err
	}
	if id == "" {
		return 0, errors.New("invalid access_token")
	}
	return strconv.ParseInt(id, 10, 64)
}

func UidInt(ctx iris.Context) int64 {
	id, err := GetUid(ctx)
	if err != nil {
		return 0
	}
	return id
}

func UidStr(ctx iris.Context) string {
	id, err := ctx.User().GetID()
	if err != nil {
		return ""
	}
	return id
}
