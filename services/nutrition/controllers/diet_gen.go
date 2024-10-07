package controllers

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func dietGenProxy(c *gin.Context) {
	remote, _ := url.Parse("http://localhost:8000")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(c.Writer, c.Request)
}