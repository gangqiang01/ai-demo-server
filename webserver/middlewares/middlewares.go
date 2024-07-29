package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/edgehook/ithings/webserver/api/jwt"
	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session, Content-Type, accesstoken, timeout, Srptoken")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "3600")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
			c.Abort()
			return
		}

		defer func() {
			if err := recover(); err != nil {
				klog.Errorf("WebServer error occurred at:%s", debug.Stack())
				responce.FailWithMessage("Server error", c)
			}
		}()

		c.Next()
	}
}

func verifyToken(token string) bool {
	if _, err := jwt.ParseCliamsToken(token); err != nil {
		return false
	}
	return true
}
