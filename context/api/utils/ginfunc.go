package utils

import (
	"appengine"
	"github.com/gamescores/gin"
	"os"
	"strings"
	"fmt"
)

const (
	gaeRootCtxKey = "GaeRootCtxKey"
	gaeCtxKey     = "GaeCtxKey"
	namespaceKey  = "Namespace"
	NamespaceHeader = "GameScoresNamespace"
)

func ResolveGaeContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeRootCtx := appengine.NewContext(c.Request)
		c.Set(gaeRootCtxKey, gaeRootCtx)

		namespace := ""

		if productionDomain := os.Getenv("PRODUCTION_DOMAIN"); productionDomain != "" {
			if strings.HasPrefix(productionDomain, ".") == false {
				productionDomain = fmt.Sprintf(".%s", productionDomain)
			}

			lastIndex := strings.LastIndex(c.Request.Host, productionDomain)

			if lastIndex > -1 {
				namespace = strings.Replace(c.Request.Host, productionDomain, "", lastIndex)
			}
		} else if devNamespace := os.Getenv("DEV_NAMESPACE"); devNamespace != "" {
			namespace = devNamespace
		}

		// Still no namespace. Maybe the request is to the appspot domain
		if namespace == "" {
			requestHost := convertDots(c.Request.Host)
			requestHost = strings.Replace(requestHost, "master.", "", 1)
			hostName, _ := appengine.ModuleHostname(gaeRootCtx, appengine.ModuleName(gaeRootCtx), "master", "")
			hostName = convertDots(hostName)
			hostName = strings.Replace(hostName, "master.", "", 1)
			hostName = fmt.Sprintf(".%s", hostName)

			lastIndex := strings.LastIndex(requestHost, hostName)

			if lastIndex > -1 {
				namespace = strings.Replace(requestHost, hostName, "", lastIndex)
			}
		}

		// Still no namespace? Last resort is a custom header
		if namespace == "" {
			namespace = c.Request.Header.Get(NamespaceHeader)
		}

		gaeRootCtx.Debugf("Using namespace: \"%s\"", namespace)
		nameSpacedGaeCtx, err := appengine.Namespace(gaeRootCtx, namespace)
		if err != nil {
			GetGaeRootContext(c).Errorf("Error creating namespace: %v", err)
			c.AbortWithError(500, err)
			return
		}

		c.Set(gaeCtxKey, nameSpacedGaeCtx)
		c.Set(namespaceKey, namespace)
	}
}

func convertDots(hostName string) string {
	return strings.Replace(hostName, "-dot-", ".", -1)
}

func GetGaeContext(c *gin.Context) appengine.Context {
	gc := c.MustGet(gaeCtxKey)
	return gc.(appengine.Context)
}

func GetGaeRootContext(c *gin.Context) appengine.Context {
	gc := c.MustGet(gaeRootCtxKey)
	return gc.(appengine.Context)
}

func GetNamespace(c *gin.Context) string {
	gc := c.MustGet(namespaceKey)
	return gc.(string)
}

func AbortWithError(c *gin.Context, err error) {
	GetGaeRootContext(c).Errorf("Error: %v", err)
	c.AbortWithError(500, err)
}
