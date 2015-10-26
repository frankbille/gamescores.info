package service

import (
	gin "github.com/gamescores/gin"
	"api/dao"
	"api/utils"
	"api/domain"
)

const (
	gameContextKey = "gameContext"
)

type ContextDefinitionService struct {
}

func CreateContextDefinitionService() ContextDefinitionService {
	return ContextDefinitionService{}
}

func (cds ContextDefinitionService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	parentRoute.GET("/context", cds.getContext)
}

func (cds ContextDefinitionService) getContext(c *gin.Context) {
	c.JSON(200, GetGameContext(c))
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

func GetGameContext(c *gin.Context) *domain.ContextDefinition {
	gc := c.MustGet(gameContextKey)
	return gc.(*domain.ContextDefinition)
}
