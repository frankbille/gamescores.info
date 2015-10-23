package defaultapp

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"fmt"
	gin "github.com/gamescores/gin"
	"net/http"
	"os"
	"regexp"
)

const (
	relPrepare RelType = "prepare"
	relCheckID RelType = "checkid"
)

var (
	ErrIDAlreadyExists = errors.New("There already exists a context with the ID")
)

type contextDefinitionService struct {
}

func createContextDefinitionService() contextDefinitionService {
	return contextDefinitionService{}
}

func (cds contextDefinitionService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	contextRoute := parentRoute.Group("/context")

	contextRoute.GET("/prepare", mustBeAuthenticated(), cds.prepareNewContext)
	contextRoute.GET("/checkid/:id", mustBeAuthenticated(), cds.checkID)
	contextRoute.POST("", mustBeAuthenticated(), cds.saveContext)
}

func (cds contextDefinitionService) prepareNewContext(c *gin.Context) {
	newContextDefinition := ContextDefinition{}
	newContextDefinition.Active = true

	newContextDefinition.AddLink(relCheckID, "/api/context/checkid/")
	newContextDefinition.AddLink(relCreate, "/api/context")

	c.JSON(http.StatusOK, newContextDefinition)
}

func (cds contextDefinitionService) checkID(c *gin.Context) {
	var checkResult struct {
		ID          string `json:"id"`
		IDAvailable bool   `json:"available"`
		Valid       bool   `json:"valid"`
	}

	checkResult.ID = c.Params.ByName("id")

	if cds._isIDValid(checkResult.ID) {
		contextDefinitionDao := createContextDefinitionDao(c)
		idExists, err := contextDefinitionDao.checkIDExists(checkResult.ID)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		checkResult.IDAvailable = idExists == false
		checkResult.Valid = true
	} else {
		checkResult.Valid = false
	}

	c.JSON(http.StatusOK, checkResult)

}

func (cds contextDefinitionService) saveContext(c *gin.Context) {
	var contextDefinition ContextDefinition

	c.Bind(&contextDefinition)

	if !cds._isValid(contextDefinition) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := getCurrentUserFromGinContext(c)

	// Run all of this in a transaction (even if we are on service level)
	err := datastore.RunInTransaction(getGaeContext(c), func(gaeCtx appengine.Context) error {
		contextDao := contextDefinitionDao{dao{gaeCtx}}
		idExists, err := contextDao.checkIDExists(contextDefinition.ID)

		if err != nil {
			return err
		}

		if idExists {
			return ErrIDAlreadyExists
		}

		contextDefinition.Owner = datastore.NewKey(gaeCtx, entityUser, user.UserID, 0, nil)
		contextDefinition.Active = true

		err = contextDao.saveContext(contextDefinition)

		if err != nil {
			return err
		}

		return nil
	}, &datastore.TransactionOptions{
		XG: true,
	})

	if err != nil {
		returnStatus := http.StatusInternalServerError
		if err == ErrIDAlreadyExists {
			returnStatus = http.StatusConflict
		}

		c.AbortWithError(returnStatus, err)
		return
	}

	contextDefinition.RemoveLink(relCreate)
	contextDefinition.RemoveLink(relCheckID)

	if productionDomain := os.Getenv("PRODUCTION_DOMAIN"); productionDomain != "" {
		contextDefinition.AddLink(relSelf, fmt.Sprintf("http://%s.%s", contextDefinition.ID, productionDomain))
	}

	c.JSON(200, contextDefinition)
}

func (cds contextDefinitionService) _isValid(contextDefinition ContextDefinition) bool {
	return cds._isIDValid(contextDefinition.ID)
}

func (cds contextDefinitionService) _isIDValid(ID string) bool {
	re := regexp.MustCompile("^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$")
	return re.MatchString(ID)
}
