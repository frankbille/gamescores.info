package context

import (
	gin "github.com/gamescores/gin"
)

type adminService struct {
}

func createAdminService() adminService {
	return adminService{}
}

func (as adminService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	adminRoute := parentRoute.Group("/admin")

	importRoute := adminRoute.Group("/import")
	importRoute.GET("/preparescoreboardv1", mustBeAdmin(), as.prepareImportScoreBoardV1)
	importRoute.POST("/scoreboardv1", mustBeAdmin(), as.importScoreBoardV1)
	importRoute.GET("/scoreboardv1/status", mustBeAdmin(), as.importScoreBoardV1Status)

	tasksRoute := rootRoute.Group("/tasks")
	tasksRoute.POST("/import/scoreboardv1", as.doImportScoreBoardV1)
}
