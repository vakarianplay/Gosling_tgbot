package main

import (
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// var dbTick *sql.DB
// var contentTbl string

func Ticker(port, baud, contentCmd, authorCmd string) {

	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		log.Println("tick "+port, " - "+baud)
		ans, au := viewBaseLine()
		ans = contentCmd + "\"" + ans + "\""
		au = authorCmd + "\"" + au + "\""
		log.Println(ans, "    ", au)
	}
}

func RunTick() {
	// go Ticker()
	// select {}
}

func viewBaseLine() (string, string) {
	qRnd := strings.Replace(QgetRandom, "{table}", contentTable, -1)
	result, err := db.Query(qRnd)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var id int
	var baseLine string
	var author string
	var tgid string
	result.Next()
	if err := result.Scan(&id, &baseLine, &author, &tgid); err != nil {
		panic(err)
	}
	ans := cleanString(baseLine)
	// log.Println(ans)
	return ans, author
}

func cleanString(input string) string {
	cleaned := strings.ReplaceAll(input, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\"", "")
	cleaned = strings.ReplaceAll(cleaned, "'", "")

	return cleaned
}
