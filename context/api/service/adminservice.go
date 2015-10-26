package service

import (
	gin "github.com/gamescores/gin"
)

type AdminService struct {
}

func CreateAdminService() AdminService {
	return AdminService{}
}

func (as AdminService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	adminRoute := parentRoute.Group("/admin")

	importRoute := adminRoute.Group("/import")
	importRoute.GET("/preparescoreboardv1", mustBeAdmin(), as.prepareImportScoreBoardV1)
	importRoute.POST("/scoreboardv1", mustBeAdmin(), as.importScoreBoardV1)
	importRoute.GET("/scoreboardv1/status", mustBeAdmin(), as.importScoreBoardV1Status)

	tasksRoute := rootRoute.Group("/tasks")
	tasksRoute.POST("/import/scoreboardv1", as.doImportScoreBoardV1)
}
