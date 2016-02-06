package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/kardianos/osext"
	_ "github.com/mattn/go-sqlite3"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	Path     string
	logger   log.Logger
	serverId int
)

const (
	TableName = "bandwidth"
)

func init() {
	logger = log.New("module", "cli.main")
	flag.IntVar(&serverId, "server", -1, "server code. Use netspeed-cli --list to retrieve correct code")
}

// If there's a problem to connect, this is the output:
/*
	Retrieving speedtest.net configuration...
	Retrieving speedtest.net server list...
	Failed to parse list of speedtest.net servers
*/

// Execute receives the full string with the command to run and retrieve the internet info.
func execute(command string) ([]byte, error) {
	splitted := strings.Split(command, " ")
	head := splitted[0]
	parts := splitted[1:]

	logger.Debug("Parsing command to retrieve speed status", "command", head, "parameters", strings.Join(parts, " "))
	return exec.Command(head, parts...).Output()
}

// AddPoint will insert a new row in database that represents the instant status of the internet connection.
func addPoint(id int, ping, download, upload float32) (int64, error) {
	timestamp := time.Now()
	logger.Info("Opening database", "path", Path)
	db, err := sql.Open("sqlite3", Path)
	defer db.Close()
	if err != nil {
		logger.Crit("Error opening database. Program terminated.", "error", err)
		panic(err)
	}

	qCreate := `CREATE TABLE IF NOT EXISTS %s(
		id INTEGER NOT NULL PRIMARY KEY, 
		timestamp TEXT, 
		ping REAL, 
		download REAL, 
		upload REAL, 
		serverid TEXT)`
	_, err = db.Exec(fmt.Sprintf(qCreate, TableName))
	if err != nil {
		logger.Error("Error while creating the table (if it does not already exists). Proceeding...", "error", err)
	}

	qInsert := `INSERT INTO %s(
		timestamp, 
		ping, 
		download, 
		upload, 
		serverid)
		values(?,?,?,?,?)`
	stmt, err := db.Prepare(fmt.Sprintf(qInsert, TableName))
	defer stmt.Close()
	if err != nil {
		logger.Crit("Failed to prepare query to insert the new row into datatable. Program terminated!",
			"error", err,
			"timestamp", timestamp,
			"ping", ping,
			"download", download,
			"upload", upload,
			"id", id)
		panic(err)
	}

	res, err := stmt.Exec(timestamp, ping, download, upload, id)
	if err != nil {
		logger.Crit("Failed to insert the new row into datatable. Program terminated!",
			"error", err,
			"timestamp", timestamp,
			"ping", ping,
			"download", download,
			"upload", upload,
			"id", id)
		panic(err)
	}
	return res.LastInsertId()
}

func main() {
	//serverId := "5151"
	flag.Parse()

	if serverId == -1 {
		logger.Crit("Server ID has not been defined. Program terminated.")
		os.Exit(1)
	}

	var err error
	Path, err = osext.ExecutableFolder()
	Path += "/values.db"
	if err != nil {
		logger.Error("Error retrieving program folder location. Maybe it doesn't exist...")
	}

	command := fmt.Sprintf("speedtest-cli --simple --server %d", serverId)
	out, err := execute(command)
	if err != nil {
		logger.Error("Output from program execution", "server", serverId)
		panic(err)
	}

	splittedOut := strings.Split(string(out), " ")
	ping := splittedOut[1]
	download := splittedOut[3]
	upload := splittedOut[5]

	parsedPing, err := strconv.ParseFloat(ping, 32)
	if err != nil {
		logger.Error("Error parsing ping value", "ping", ping)
	}
	parsedDL, err := strconv.ParseFloat(download, 32)
	if err != nil {
		logger.Error("Failed to convert download value.", "", download)
	}

	parsedUL, err := strconv.ParseFloat(upload, 32)
	if err != nil {
		logger.Error("Error parsing uoload value", "ping", upload)
	}

	logger.Debug("Data collected", "ping", parsedPing, "download", parsedDL, "upload", parsedUL)
	addPoint(serverId, float32(parsedPing), float32(parsedDL), float32(parsedUL))
	//fmt.Printf("Ping: %s\nDownload: %s\nUpload: %s\nDone.\n", parsedPing, parsedDL, parsedUL)
	logger.Debug("Program finished")
}
