package rest

import (
	"appengine"
	"fmt"
	gin "github.com/gamescores/gin"
	http "net/http"
	"os"
	"strings"
	"src/service"
	"src/utils"
)

func init() {
	r := gin.New()

	root := r.Group("/")

	root.Use(gaeContext())
	root.Use(service.ResolveGameContext())
	root.Use(service.ResolveUser())

	api := root.Group("/api")

	// Create list of services used
	services := []service.RestService{
		service.CreateContextDefinitionService(),
		service.CreateUserService(),
		service.CreatePlayerService(),
		service.CreateLeagueService(),
		service.CreateGameService(),
		service.CreateAdminService(),
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
		c.Set(utils.GaeRootCtxKey, gaeRootCtx)

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
			namespace = c.Request.Header.Get(service.NamespaceHeader)
		}

		gaeRootCtx.Debugf("Using namespace: \"%s\"", namespace)
		nameSpacedGaeCtx, err := appengine.Namespace(gaeRootCtx, namespace)
		if err != nil {
			utils.GetGaeRootContext(c).Errorf("Error creating namespace: %v", err)
			c.AbortWithError(500, err)
			return
		}

		c.Set(utils.GaeCtxKey, nameSpacedGaeCtx)
		c.Set(utils.NamespaceKey, namespace)
	}
}

func convertDots(hostName string) string {
	return strings.Replace(hostName, "-dot-", ".", -1)
}
