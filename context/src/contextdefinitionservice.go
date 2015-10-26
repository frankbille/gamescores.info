package context

import (
	gin "github.com/gamescores/gin"
)

const (
	gameContextKey = "gameContext"
)

type contextDefinitionService struct {
}

func createContextDefinitionService() contextDefinitionService {
	return contextDefinitionService{}
}

func (cds contextDefinitionService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	parentRoute.GET("/context", cds.getContext)
}

func (cds contextDefinitionService) getContext(c *gin.Context) {
	c.JSON(200, getGameContext(c))
}

func resolveGameContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		contextDefinitionDao := createContextDefinitionDao(c)

		contextDefinition, err := contextDefinitionDao.getContext(getNamespace(c))

		if err != nil {
			abortWithError(c, err)
			return
		}

		c.Set(gameContextKey, contextDefinition)
	}
}

func getGameContext(c *gin.Context) *ContextDefinition {
	gc := c.MustGet(gameContextKey)
	return gc.(*ContextDefinition)
}
