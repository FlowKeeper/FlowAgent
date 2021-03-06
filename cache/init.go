package cache

import (
	"database/sql"
	"log"

	"github.com/FlowKeeper/FlowAgent/v2/config"
	"github.com/google/uuid"

	//Black import because we are using sqlite here to cache item results
	_ "github.com/mattn/go-sqlite3"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
)

var db *sql.DB

const loggingArea = "CACHE"

//Init initializes the sqlite database and its client
func Init() {
	var err error
	logger.Info(loggingArea, "Initializing cache")

	db, err = sql.Open("sqlite3", config.Config.CachePath)
	if err != nil {
		logger.Fatal(loggingArea)
	}

	sqlStmt := "CREATE TABLE IF NOT EXISTS config (id INTEGER PRIMARY KEY NOT NULL, scraperID TEXT, agentID TEXT NOT NULL, version INTEGER NOT NULL);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(loggingArea, err)
	}

	sqlStmt = `CREATE TABLE IF NOT EXISTS results (
		ItemID TEXT NOT NULL,
		CapturedAt TEXT NOT NULL,
		ValueString TEXT,
		ValueInt TEXT,
		Type INTEGER NOT NULL,
		Error TEXT
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(loggingArea, err)
	}

	versionCheck()
	logger.Info(loggingArea, "Cache initialization complete")
}

const currentSQLVersion = 1

func versionCheck() {
	res, err := db.Query(`SELECT scraperID,agentID,version FROM config`)
	if err != nil {
		log.Fatal(loggingArea, err)
	}

	if !res.Next() {
		initializeConfig()
		return
	}

	var version int
	if err := res.Scan(&Config.ScraperUUID, &Config.AgentUUID, &version); err != nil {
		logger.Fatal(loggingArea, err)
	}

	res.Close()
	if version == currentSQLVersion {
		logger.Info(loggingArea, "Local Agent Configuration is right verison. No need to upgrade.")
		return
	}

	logger.Fatal(loggingArea, "NOT IMPLEMENTED: SQlite Configuration needed upgrade.")
	//ToDo: Upgrade procedures for the sqllite if needed
}

func initializeConfig() {
	logger.Info(loggingArea, "Initializing local agent config")
	if _, err := db.Exec(`INSERT INTO config (id, agentID, version) VALUES (?,?,?)`, 0, uuid.New(), currentSQLVersion); err != nil {
		logger.Fatal(loggingArea, err)
	}
}
