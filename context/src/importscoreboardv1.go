package context

import (
	"appengine"
	"appengine/datastore"
	"appengine/taskqueue"
	"appengine/urlfetch"
	"encoding/xml"
	"github.com/gamescores/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	relImport       RelType = "import"
	relImportStatus RelType = "importstatus"
)

type ScoreBoardV1Import struct {
	DefaultHalResource
	DbDumpUrl string `json:"dbDumpUrl"`
}

type ScoreBoardV1ImportStatus struct {
	DefaultHalResource
	Importing           bool `json:"importing"`
	ImportedPlayerCount int  `json:"importedPlayerCount"`
	TotalPlayerCount    int  `json:"totalPlayerCount"`
	ImportedLeagueCount int  `json:"importedLeagueCount"`
	TotalLeagueCount    int  `json:"totalLeagueCount"`
	ImportedGameCount   int  `json:"importedGameCount"`
	TotalGameCount      int  `json:"totalGameCount"`
}

func (as adminService) prepareImportScoreBoardV1(c *gin.Context) {
	importDefinition := ScoreBoardV1Import{}

	importDefinition.AddLink(relImport, "/api/admin/import/scoreboardv1")

	c.JSON(200, importDefinition)
}

func (as adminService) importScoreBoardV1(c *gin.Context) {
	var importDefinition ScoreBoardV1Import

	c.Bind(&importDefinition)

	if importDefinition.DbDumpUrl == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	createTask := &taskqueue.Task{
		Path:    "/tasks/import/scoreboardv1",
		Payload: []byte(importDefinition.DbDumpUrl),
	}
	hostName, _ := appengine.ModuleHostname(getGaeRootContext(c), appengine.ModuleName(getGaeRootContext(c)), "", "")
	createTask.Header = http.Header{}
	createTask.Header.Set("Host", hostName)
	createTask.Header.Set(namespaceHeader, getNamespace(c))

	_, err := taskqueue.Add(getGaeRootContext(c), createTask, "contextqueue")

	if err != nil {
		getGaeRootContext(c).Errorf("Error calling taskqueue.Add in importScoreBoardV1: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	importDao := createImportDao(c)
	importStatus, err := importDao.setStatus(true, 0, 0, 0, 0, 0, 0)

	if err != nil {
		getGaeRootContext(c).Errorf("Error calling importDao.setStatus in importScoreBoardV1: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	importStatus.AddLink(relImportStatus, "/api/admin/import/scoreboardv1/status")

	c.JSON(200, importStatus)
}

func (as adminService) importScoreBoardV1Status(c *gin.Context) {
	importDao := createImportDao(c)

	importStatus, err := importDao.getStatus()

	if err != nil {
		getGaeRootContext(c).Errorf("%v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if importStatus == nil {
		importStatus = &ScoreBoardV1ImportStatus{}
	}

	importStatus.AddLink(relImportStatus, "/api/admin/import/scoreboardv1/status")

	c.JSON(200, importStatus)
}

func (as adminService) doImportScoreBoardV1(c *gin.Context) {
	getGaeRootContext(c).Infof("%#v", c.Request)

	importDao := createImportDao(c)

	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		getGaeRootContext(c).Errorf("Error calling ioutil.ReadAll(c.Request.Body) in doImportScoreBoardV1: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		importDao.setStatus(false, 0, 0, 0, 0, 0, 0)
		return
	}

	dbDumpUrl := string(body)

	httpClient := urlfetch.Client(getGaeRootContext(c))
	response, err := httpClient.Get(dbDumpUrl)

	if err != nil {
		getGaeRootContext(c).Errorf("Error calling httpClient.Get in doImportScoreBoardV1: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		importDao.setStatus(false, 0, 0, 0, 0, 0, 0)
		return
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		getGaeRootContext(c).Errorf("Error calling ioutil.ReadAll(response.Body) in doImportScoreBoardV1: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		importDao.setStatus(false, 0, 0, 0, 0, 0, 0)
		return
	}

	dump := MysqlDump{}

	err = xml.Unmarshal(data, &dump)

	if err != nil {
		getGaeRootContext(c).Errorf("Error calling xml.Unmarshal in doImportScoreBoardV1: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		importDao.setStatus(false, 0, 0, 0, 0, 0, 0)
		return
	}

	database := dump.Databases[0]

	playerTable := as._getTableByName(database, "player")
	leagueTable := as._getTableByName(database, "league")
	gameTable := as._getTableByName(database, "game")
	gameTeamTable := as._createLookupMap(as._getTableByName(database, "game_team"), "id")
	teamPlayersTable := as._createLookupMap(as._getTableByName(database, "team_players"), "team_id")

	playerTotal := len(playerTable.Rows)
	playerCount := 0
	leagueTotal := len(leagueTable.Rows)
	leagueCount := 0
	gameTotal := len(gameTable.Rows)
	gameCount := 0
	_, err = importDao.setStatus(true, playerTotal, playerCount, leagueTotal, leagueCount, gameTotal, gameCount)

	if err != nil {
		getGaeRootContext(c).Errorf("importDao.setStatus: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Add players first
	as._deleteAll(entityPlayer, getGaeContext(c))
	playerDao := createPlayerDao(c)

	playerConvertIdMap := make(map[string]int64)
	for _, playerRow := range playerTable.Rows {
		id := as._getFieldValueByName(playerRow, "id")
		name := as._getFieldValueByName(playerRow, "name")

		savedPlayer, err := playerDao.savePlayer(Player{
			Active: true,
			Name:   name,
		})

		if err != nil {
			getGaeRootContext(c).Errorf("Error calling playerDao.savePlayer in doImportScoreBoardV1: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			importDao.setStatus(false, 0, 0, 0, 0, 0, 0)
			return
		}

		playerConvertIdMap[id] = savedPlayer.ID

		playerCount++
		_, err = importDao.setStatus(true, playerTotal, playerCount, leagueTotal, leagueCount, gameTotal, gameCount)

		if err != nil {
			getGaeRootContext(c).Errorf("importDao.setStatus: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Add leagues
	as._deleteAll(entityLeague, getGaeContext(c))
	leagueDao := createLeagueDao(c)

	leagueConvertIdMap := make(map[string]int64)
	for _, leagueRow := range leagueTable.Rows {
		id := as._getFieldValueByName(leagueRow, "id")
		name := as._getFieldValueByName(leagueRow, "name")

		savedLeague, err := leagueDao.saveLeague(League{
			Active: true,
			Name:   name,
		})

		if err != nil {
			getGaeRootContext(c).Errorf("Error calling leagueDao.saveLeague in doImportScoreBoardV1: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			importDao.setStatus(false, 0, 0, 0, 0, 0, 0)
			return
		}

		leagueConvertIdMap[id] = savedLeague.ID

		leagueCount++
		_, err = importDao.setStatus(true, playerTotal, playerCount, leagueTotal, leagueCount, gameTotal, gameCount)

		if err != nil {
			getGaeRootContext(c).Errorf("importDao.setStatus: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Add games
	as._deleteAll(entityGame, getGaeContext(c))
	gameDao := createGameDao(c)

	for _, gameRow := range gameTable.Rows {
		gameDate := as._getFieldDateValueByName(gameRow, "game_date")
		team1IDString := as._getFieldValueByName(gameRow, "team1_id")
		team2IDString := as._getFieldValueByName(gameRow, "team2_id")
		leagueIDString := as._getFieldValueByName(gameRow, "league_id")

		team1 := as._createTeam(team1IDString, gameTeamTable, teamPlayersTable, playerConvertIdMap)
		team2 := as._createTeam(team2IDString, gameTeamTable, teamPlayersTable, playerConvertIdMap)

		game := Game{
			GameDate: gameDate,
			LeagueID: leagueConvertIdMap[leagueIDString],
			Team1:    team1,
			Team2:    team2,
		}

		_, err := gameDao.saveGame(game)

		if err != nil {
			getGaeRootContext(c).Errorf("Error calling gameDao.saveGame in doImportScoreBoardV1: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			importDao.setStatus(false, 0, 0, 0, 0, 0, 0)
			return
		}

		gameCount++
		_, err = importDao.setStatus(true, playerTotal, playerCount, leagueTotal, leagueCount, gameTotal, gameCount)

		if err != nil {
			getGaeRootContext(c).Errorf("importDao.setStatus: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	_, err = importDao.setStatus(false, playerTotal, playerCount, leagueTotal, leagueCount, gameTotal, gameCount)

	if err != nil {
		getGaeRootContext(c).Errorf("importDao.setStatus: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (as adminService) _createTeam(idString string, gameTeamTable, teamPlayersTable map[string][]Row, playerConvertIdMap map[string]int64) GameTeam {
	gameTeamRow := gameTeamTable[idString][0]

	teamIdString := as._getFieldValueByName(gameTeamRow, "team_id")
	score := as._getFieldIntValueByName(gameTeamRow, "score")

	playerRows := teamPlayersTable[teamIdString]
	playerIds := make([]int64, len(playerRows))
	for idx, teamPlayersRow := range playerRows {
		playerIds[idx] = playerConvertIdMap[as._getFieldValueByName(teamPlayersRow, "player_id")]
	}

	return GameTeam{
		Score:   score,
		Players: playerIds,
	}
}

func (as adminService) _createLookupMap(table *TableData, keyFieldName string) map[string][]Row {
	lookupMap := make(map[string][]Row)
	for _, row := range table.Rows {
		keyFieldValue := as._getFieldValueByName(row, keyFieldName)
		rows, found := lookupMap[keyFieldValue]
		if !found {
			rows = []Row{}
		}
		rows = append(rows, row)
		lookupMap[keyFieldValue] = rows
	}
	return lookupMap
}

func (as adminService) _getFieldIntValueByName(row Row, fieldName string) int {
	valueString := as._getFieldValueByName(row, fieldName)
	if valueString != "" {
		value, _ := strconv.ParseInt(valueString, 10, 0)
		return int(value)
	}
	return 0
}

func (as adminService) _getFieldDateValueByName(row Row, fieldName string) time.Time {
	valueString := as._getFieldValueByName(row, fieldName)
	if valueString != "" {
		value, _ := time.Parse("2006-01-02", valueString)
		return value
	}
	return time.Now()
}

func (as adminService) _getFieldValueByName(row Row, fieldName string) string {
	field := as._getFieldByName(row, fieldName)
	if field != nil {
		return field.Value
	} else {
		return ""
	}
}

func (as adminService) _getFieldByName(row Row, fieldName string) *Field {
	for _, field := range row.Fields {
		if field.Name == fieldName {
			return &field
		}
	}
	return nil
}

func (as adminService) _getTableByName(database Database, tableName string) *TableData {
	for _, table := range database.Tables {
		if table.Name == tableName {
			return &table
		}
	}
	return nil
}

func (as adminService) _deleteAll(kind string, c appengine.Context) {
	q := datastore.NewQuery(kind)
	q = q.KeysOnly()
	keys, _ := q.GetAll(c, nil)
	datastore.DeleteMulti(c, keys)
}

type Field struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",innerxml"`
}

type Row struct {
	Fields []Field `xml:"field"`
}

type TableData struct {
	Name string `xml:"name,attr"`
	Rows []Row  `xml:"row"`
}

type Database struct {
	Tables []TableData `xml:"table_data"`
}

type MysqlDump struct {
	Databases []Database `xml:"database"`
}
