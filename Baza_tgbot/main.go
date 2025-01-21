package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	yaml "gopkg.in/yaml.v3"
)

func main() {

	botapi, dbPath, contentDB, usersDB := readCfg()[0], readCfg()[1], readCfg()[2], readCfg()[3]
	srvPath, port := readCfg()[4], readCfg()[5]
	infFlag, _ := strconv.ParseBool(readCfg()[6])
	itteraction, _ := strconv.Atoi(readCfg()[7])
	infoTemplate := readCfg()[8]
	infoUrls := readCfg()[9]

	log.Println(botapi, dbPath, contentDB, usersDB)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(dbPath + " connected")
	}

	startInformer(infFlag, infoTemplate, infoUrls)
	go HTTPServer(srvPath, port, itteraction)
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

	info_flag := (cfgYaml["informer"].(map[string]interface{})["enabled"])
	itter := (cfgYaml["informer"].(map[string]interface{})["itteraction"])
	info_tmp := (cfgYaml["informer"].(map[string]interface{})["info_str"])
	url_list := (cfgYaml["informer"].(map[string]interface{})["urls"])

	apiKey_ := fmt.Sprintf("%v", apiKey)
	dbPath_ := fmt.Sprintf("%v", dbPath)
	contentDB_ := fmt.Sprintf("%v", contentDB)
	usersDB_ := fmt.Sprintf("%v", usersDB)

	root_path_ := fmt.Sprintf("%v", root_serv)
	port_s := fmt.Sprintf("%v", port_srv)

	info_flag_ := fmt.Sprintf("%v", info_flag)
	itter_ := fmt.Sprintf("%v", itter)
	info_str_ := fmt.Sprintf("%v", info_tmp)
	urls_ := fmt.Sprintf("%v", url_list)

	var out []string
	out = append(out, apiKey_, dbPath_, contentDB_, usersDB_, root_path_, port_s, info_flag_, itter_, info_str_, urls_)

	return out
}
