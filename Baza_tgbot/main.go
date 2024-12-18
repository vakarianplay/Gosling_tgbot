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

	go Ticker(readCfg()[4], readCfg()[5])
	TelegramBot(botapi, contentDB, usersDB, db)

}

func readCfg() []string {

	var cfgYaml map[string]interface{}
	// cfgFile, err := os.ReadFile("config.yaml")
	cfgFile, err := ioutil.ReadFile("config.yml")
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

	root_serv := (cfgYaml["server_config"].(map[string]interface{})["root_path"])
	port_srv := (cfgYaml["server_config"].(map[string]interface{})["port"])

	apiKey_ := fmt.Sprintf("%v", apiKey)
	dbPath_ := fmt.Sprintf("%v", dbPath)
	contentDB_ := fmt.Sprintf("%v", contentDB)
	usersDB_ := fmt.Sprintf("%v", usersDB)

	root_path_ := fmt.Sprintf("%v", root_serv)
	port_s := fmt.Sprintf("%v", port_srv)

	var out []string
	out = append(out, apiKey_, dbPath_, contentDB_, usersDB_, root_path_, port_s)

	return out
}
