package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var InfoGl string
var counter int
var response string

func HTTPServer(root_path, port_ string, itter_ int) {
	http.HandleFunc(root_path, func(w http.ResponseWriter, r *http.Request) {
		log.Println(itter_, "  -----  ", counter)
		if itter_ > counter {
			ans, author := viewBaseLine()
			response = fmt.Sprintf("%s   Автор: %s", ans, author)
			counter++
		} else {
			now := time.Now()
			currentTime := now.Format("02-01 15:04")
			response = currentTime + "              " + InfoGl
			counter = 0
		}

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

func startInformer(enabled_ bool, infTemplate_ string, infUrls_ string) {
	if enabled_ {
		InfoGl = informer(enabled_, infTemplate_, infUrls_)
		log.Println(InfoGl)
		ticker := time.NewTicker(30 * time.Minute)
		go func() {
			for {
				<-ticker.C
				InfoGl = informer(enabled_, infTemplate_, infUrls_)
				log.Println("update: " + InfoGl)
			}
		}()
	}
}

func informer(enabled_ bool, infTemplate_ string, infUrls_ string) string {
	// now := time.Now()
	// currentTime := now.Format("02-01 15:04")

	// infoAns := currentTime + "         "
	var infoStr string
	if enabled_ {
		result, err := replacePlaceholders(infTemplate_, infUrls_)
		if err != nil {
			log.Println("Error:", err)
		} else {
			infoStr = result
		}
	}
	return infoStr
}

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

func replacePlaceholders(template string, urls string) (string, error) {
	urls = strings.ReplaceAll(urls, " ", "")
	urlList := strings.Split(urls, ",")
	placeholderCount := strings.Count(template, "{placeholder}")

	if placeholderCount != len(urlList) && len(urls) > 0 {
		return "", fmt.Errorf("placeholders (%d) dont't eq URL (%d)", placeholderCount, len(urlList))
	}
	var sb strings.Builder
	start := 0
	urlIndex := 0

	for {
		placeholderStart := strings.Index(template[start:], "{placeholder}")
		if placeholderStart == -1 {
			sb.WriteString(template[start:])
			break
		}

		sb.WriteString(template[start : start+placeholderStart])

		if urlIndex < len(urlList) {
			sb.WriteString(strings.TrimSpace(getInfo(urlList[urlIndex])))
			urlIndex++
		}

		start += placeholderStart + len("{placeholder}")
	}
	return sb.String(), nil
}
