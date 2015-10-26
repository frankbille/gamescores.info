package dao

import (
	"api/domain"
	"api/utils"
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
)

const EntityImport string = "Import"

type ImportDao struct {
	dao
}

func CreateImportDao(c *gin.Context) ImportDao {
	dao := createDao(utils.GetGaeContext(c))
	return ImportDao{dao}
}

func (dao *ImportDao) GetStatus() (*domain.ScoreBoardV1ImportStatus, error) {
	var importStatus domain.ScoreBoardV1ImportStatus

	key := datastore.NewKey(dao.Context, EntityImport, "", 1, nil)

	err := dao.get(key, &importStatus)

	if err == datastore.ErrNoSuchEntity {
		return nil, nil
	}

	return &importStatus, err
}

func (dao *ImportDao) SetStatus(importing bool, playerTotal, playerCreated, leagueTotal, leagueCreated, gameTotal, gameCreated int) (*domain.ScoreBoardV1ImportStatus, error) {
	importStatus, err := dao.GetStatus()

	if err != nil {
		return nil, err
	}

	if importStatus == nil {
		importStatus = &domain.ScoreBoardV1ImportStatus{}
	}

	importStatus.Importing = importing
	importStatus.TotalPlayerCount = playerTotal
	importStatus.ImportedPlayerCount = playerCreated
	importStatus.TotalLeagueCount = leagueTotal
	importStatus.ImportedLeagueCount = leagueCreated
	importStatus.TotalGameCount = gameTotal
	importStatus.ImportedGameCount = gameCreated

	key := datastore.NewKey(dao.Context, EntityImport, "", 1, nil)

	_, err = dao.save(key, importStatus)

	return importStatus, err
}
