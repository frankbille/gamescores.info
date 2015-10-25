package utils

import (
	"appengine"
	"github.com/gamescores/gin"
)

const (
	GaeRootCtxKey = "GaeRootCtxKey"
	GaeCtxKey     = "GaeCtxKey"
	NamespaceKey  = "Namespace"
)

func GetGaeContext(c *gin.Context) appengine.Context {
	gc := c.MustGet(GaeCtxKey)
	return gc.(appengine.Context)
}

func GetGaeRootContext(c *gin.Context) appengine.Context {
	gc := c.MustGet(GaeRootCtxKey)
	return gc.(appengine.Context)
}

func GetNamespace(c *gin.Context) string {
	gc := c.MustGet(NamespaceKey)
	return gc.(string)
}

func AbortWithError(c *gin.Context, err error) {
	GetGaeRootContext(c).Errorf("Error: %v", err)
	c.AbortWithError(500, err)
}
