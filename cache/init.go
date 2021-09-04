package cache

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowagent/v2/config"
)

var db *sql.DB

const loggingArea = "CACHE"

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
	if err := res.Scan(&Config.ScraperID, &Config.AgentID, &version); err != nil {
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
