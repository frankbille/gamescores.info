package service

import (
	gin "github.com/gamescores/gin"
	"src/dao"
	"src/domain"
	"src/utils"
)

const (
	gameContextKey = "gameContext"
)

type contextDefinitionService struct {
}

func CreateContextDefinitionService() contextDefinitionService {
	return contextDefinitionService{}
}

func (cds contextDefinitionService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	parentRoute.GET("/context", cds.getContext)
}

func (cds contextDefinitionService) getContext(c *gin.Context) {
	c.JSON(200, getGameContext(c))
}

func ResolveGameContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		contextDefinitionDao := dao.CreateContextDefinitionDao(c)

		contextDefinition, err := contextDefinitionDao.GetContext(utils.GetNamespace(c))

		if err != nil {
			utils.AbortWithError(c, err)
			return
		}

		c.Set(gameContextKey, contextDefinition)
	}
}

func getGameContext(c *gin.Context) *domain.ContextDefinition {
	gc := c.MustGet(gameContextKey)
	return gc.(*domain.ContextDefinition)
}
