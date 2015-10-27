package context

import (
	"api/service"
	"api/utils"
	gin "github.com/gamescores/gin"
	http "net/http"
)

type restService interface {
	CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup)
}

func init() {
	r := gin.New()

	root := r.Group("/")

	root.Use(utils.ResolveGaeContext())
	root.Use(service.ResolveGameContext())
	root.Use(service.ResolveUser())

	api := root.Group("/api")

	// Create list of services used
	services := []restService{
		service.CreateContextDefinitionService(),
		service.CreateUserService(),
		service.CreatePlayerService(),
		service.CreateLeagueService(),
		service.CreateGameService(),
		service.CreateRatingService(),
		service.CreateAdminService(),
	}

	// Process the services
	for _, service := range services {
		service.CreateRoutes(api, root)
	}

	http.Handle("/", r)
}
