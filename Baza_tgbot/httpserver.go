package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func HTTPServer(root_path, port_ string) {

	http.HandleFunc(root_path, func(w http.ResponseWriter, r *http.Request) {
		ans, author := viewBaseLine()
		response := fmt.Sprintf("%s   Автор: %s", ans, author)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
	log.Println("Server listening on ", port_)
	http.ListenAndServe(":"+port_, nil)

}

func viewBaseLine() (string, string) {
	qRnd := strings.Replace(QgetRandom, "{table}", contentTable, -1)
	// qRnd := strings.Replace(QgetRandom, "{table}", readCfg()[2], -1)
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

// func Informer(enabled_ bool, itter_ int, infTemplate_, infUrls string) string {
// 	// urlWttr := "https://ru.wttr.in/Moscow?format=%D0%A2%D0%B5%D0%BC%D0%BF%D0%B5%D1%80%D0%B0%D1%82%D1%83%D1%80%D0%B0:%20%t%20(%f),%20%C,%20%D0%B2%D0%BB%D0%B0%D0%B6%D0%BD%D0%BE%D1%81%D1%82%D1%8C:%20%h,%20%D0%B2%D0%B5%D1%82%D0%B5%D1%80:%20%w"
// 	// urlUsd := "https://rub.rate.sx/1USDT"
// 	// urlBtc := "http://rate.sx/1BTC"

// 	// now := time.Now()
// 	// currentTime := now.Format("02-01 15:04")
// 	// infoAns := currentTime + "        " + getInfo(urlWttr) + "          RUB-USD: " + getInfo(urlUsd) + "    USD-BTC: " + getInfo(urlBtc)

// 	// return infoAns
// }

func getInfo(url string) string {

	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error seng GET:", err)
	}
	defer response.Body.Close()

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(response.Body)
	if err != nil {
		log.Fatal("Read Error:", err)
	}

	ans := buffer.String()
	return ans
}
