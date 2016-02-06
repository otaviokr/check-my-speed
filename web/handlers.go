package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"database/sql"
	"github.com/kardianos/osext"
	_ "github.com/mattn/go-sqlite3"
	log "gopkg.in/inconshreveable/log15.v2"
)

// Info wraps the data expected in the report page.
type Info struct {
	MinValue  float64
	MaxValue  float64
	AvgValue  float64
	LastValue float64
	Points    [][]interface{}
}

// Speed generates the report page.
func Speed(w http.ResponseWriter, r *http.Request) {
	currentFolder, err := osext.ExecutableFolder()
	if err != nil {
		log.Error("Could not retrieve current folder. Attempting to use dot(.) instead...", "error", err)
		currentFolder = "."
	}

	files := []string{currentFolder + "/html/templates/speed.html"}
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Crit("Failed to load template", "error", err, "templates", files)
		return
	}

	data := Info{}
	getData(&data)

	err = t.Execute(w, data)
	if err != nil {
		log.Crit("Failed to apply data to template", "error", err, "data", data)
		return
	}
}

// Index is just for tests, really.
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I'm sorry, Mario, but your report is in another page...")
}

// GetData retrieves and calculates the metrics to be displayed in the report page.
func getData(data *Info) error {
	currentFolder, err := osext.ExecutableFolder()
	if err != nil {
		log.Error("Could not retrieve current folder. Attempting to use dot(.) instead...", "error", err)
		currentFolder = "."
	}

	dbFilePath := currentFolder + "/values.db"
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Crit("Failed to opend database", "error", err, "path", dbFilePath)
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT timestamp, ping, download, upload FROM bandwidth")
	if err != nil {
		log.Crit("Could not retrieve data from database", "error", err)
		return err
	}
	defer rows.Close()

	min := 10000.0 // Unfortunately, I don't see myself with a 10000Mbit/s connection anytime soon...
	max := 0.0
	counter := 0
	average := 0.0
	data.Points = [][]interface{}{{map[string]string{"type": "datetime", "label": "Time"}, "Download", "Upload"}}

	for rows.Next() {
		var timestamp string
		var ping, download, upload float64
		rows.Scan(&timestamp, &ping, &download, &upload)

		if download < min {
			min = download
		}

		if download > max {
			max = download
		}

		average += download
		counter++

		// Timestamp is presented as YYYY-MM-DD HH:MI:SS.Milli+0000
		split := strings.Split(timestamp, " ")
		dateOnly := strings.Split(string(split[0]), "-")
		timeOnly := strings.Split(string(split[1]), ".")
		timeOnly = strings.Split(string(timeOnly[0]), ":")
		axis := fmt.Sprintf("Date(%s, %s, %s, %s, %s, %s)", dateOnly[0], dateOnly[1], dateOnly[2], timeOnly[0], timeOnly[1], timeOnly[2])

		data.Points = append(data.Points, []interface{}{axis, download, upload})
	}
	data.MinValue = min
	data.MaxValue = max
	data.AvgValue = average / float64(counter)
	data.LastValue = data.Points[len(data.Points)-1][1].(float64)

	return nil
}
