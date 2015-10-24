package context

import (
	"appengine"
	"fmt"
	gin "github.com/gamescores/gin"
	http "net/http"
	"os"
	"strings"
)

const (
	gaeRootCtxKey = "GaeRootCtxKey"
	gaeCtxKey     = "GaeCtxKey"
	namespaceKey  = "Namespace"
)

type restService interface {
	CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup)
}

func init() {
	r := gin.New()

	root := r.Group("/")

	root.Use(gaeContext())
	root.Use(resolveGameContext())
	root.Use(resolveUser())

	api := root.Group("/api")

	// Create list of services used
	services := []restService{
		createContextDefinitionService(),
		createUserService(),
		createPlayerService(),
		createLeagueService(),
		createGameService(),
		createAdminService(),
	}

	// Process the services
	for _, service := range services {
		service.CreateRoutes(api, root)
	}

	http.Handle("/", r)
}

func gaeContext() gin.HandlerFunc {
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
			requestHost = strings.Replace(requestHost, "master.", ".", 1)
			gaeRootCtx.Debugf("Request host: %s", requestHost)
			hostName, _ := appengine.ModuleHostname(gaeRootCtx, appengine.ModuleName(gaeRootCtx), "master", "")
			hostName = convertDots(hostName)
			hostName = strings.Replace(hostName, "master.", ".", 1)
			gaeRootCtx.Debugf("Hostname: %s", hostName)

			lastIndex := strings.LastIndex(c.Request.Host, hostName)

			gaeRootCtx.Debugf("Last index: %d", lastIndex)

			if lastIndex > -1 {
				namespace = strings.Replace(c.Request.Host, hostName, "", lastIndex)
			}
		}

		gaeRootCtx.Debugf("Using namespace: \"%s\"", namespace)
		nameSpacedGaeCtx, err := appengine.Namespace(gaeRootCtx, namespace)
		if err != nil {
			getGaeRootContext(c).Errorf("Error creating namespace: %v", err)
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

func getGaeContext(c *gin.Context) appengine.Context {
	gc := c.MustGet(gaeCtxKey)
	return gc.(appengine.Context)
}

func getGaeRootContext(c *gin.Context) appengine.Context {
	gc := c.MustGet(gaeRootCtxKey)
	return gc.(appengine.Context)
}

func getNamespace(c *gin.Context) string {
	gc := c.MustGet(namespaceKey)
	return gc.(string)
}

func abortWithError(c *gin.Context, err error) {
	getGaeRootContext(c).Errorf("Error: %v", err)
	c.AbortWithError(500, err)
}
