package context

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
)

const entityImport string = "Import"

type importDao struct {
	dao
}

func createImportDao(c *gin.Context) importDao {
	dao := createDao(getGaeContext(c))
	return importDao{dao}
}

func (dao *importDao) getStatus() (*ScoreBoardV1ImportStatus, error) {
	var importStatus ScoreBoardV1ImportStatus

	key := datastore.NewKey(dao.Context, entityImport, "", 1, nil)

	err := dao.get(key, &importStatus)

	if err == datastore.ErrNoSuchEntity {
		return nil, nil
	}

	return &importStatus, err
}

func (dao *importDao) setStatus(importing bool, playerTotal, playerCreated, leagueTotal, leagueCreated, gameTotal, gameCreated int) (*ScoreBoardV1ImportStatus, error) {
	importStatus, err := dao.getStatus()

	if err != nil {
		return nil, err
	}

	if importStatus == nil {
		importStatus = &ScoreBoardV1ImportStatus{}
	}

	importStatus.Importing = importing
	importStatus.TotalPlayerCount = playerTotal
	importStatus.ImportedPlayerCount = playerCreated
	importStatus.TotalLeagueCount = leagueTotal
	importStatus.ImportedLeagueCount = leagueCreated
	importStatus.TotalGameCount = gameTotal
	importStatus.ImportedGameCount = gameCreated

	key := datastore.NewKey(dao.Context, entityImport, "", 1, nil)

	_, err = dao.save(key, importStatus)

	return importStatus, err
}
