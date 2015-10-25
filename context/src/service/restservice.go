package service
import "github.com/gamescores/gin"

type RestService interface {
	CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup)
}