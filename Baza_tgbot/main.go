package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	botapi, dbPath, contentDB, usersDB := readCfg()[0], readCfg()[1], readCfg()[2], readCfg()[3]
	log.Println(botapi, dbPath, contentDB, usersDB)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(dbPath + " connected")
	}

	go Ticker()
	TelegramBot(botapi, db)

	// defer db.Close()

	// db.Exec("INSERT INTO content (text, author) VALUES ('test test', 'sample author');")
	// queryStr := strings.Replace("SELECT * FROM {table}", "{table}", contentDB, -1)
	// fmt.Println(queryStr)
	// rows, err := db.Query(queryStr)
	// if err != nil {
	// panic(err)
	// }
	// defer rows.Close() // Закрываем строки результатов при выходе из функции.

	// Перебираем строки результатов.
	// for rows.Next() {
	// 	var id int
	// 	var text string
	// 	var author string
	// 	if err := rows.Scan(&id, &text, &author); err != nil {
	// 		panic(err)
	// 	}

	// 	// Выводим содержимое строки.
	// 	fmt.Printf("ID: %d, Title: %s, Body: %s\n", id, text, author)
	// }

	// Проверяем наличие ошибок при переборе строк.
	// if err := rows.Err(); err != nil {
	// 	panic(err)
	// }
}

func readCfg() []string {

	var cfgYaml map[string]interface{}
	// cfgFile, err := os.ReadFile("config.yaml")
	cfgFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(cfgFile, &cfgYaml)

	if err != nil {
		log.Fatal(err)
	}

	apiKey := (cfgYaml["bot"].(map[string]interface{})["api_key"])
	dbPath := (cfgYaml["database"].(map[string]interface{})["file"])
	contentDB := (cfgYaml["database"].(map[string]interface{})["content_table"])
	usersDB := (cfgYaml["database"].(map[string]interface{})["users_table"])
	// font := (cfgYaml["font"].(map[string]interface{})["font_file"])
	// fsize := (cfgYaml["font"].(map[string]interface{})["font_size"])

	apiKey_ := fmt.Sprintf("%v", apiKey)
	dbPath_ := fmt.Sprintf("%v", dbPath)
	contentDB_ := fmt.Sprintf("%v", contentDB)
	usersDB_ := fmt.Sprintf("%v", usersDB)
	// font_ := fmt.Sprintf("%v", font)
	// fsize_ := fmt.Sprintf("%v", fsize)

	var out []string
	out = append(out, apiKey_, dbPath_, contentDB_, usersDB_)

	return out
}
