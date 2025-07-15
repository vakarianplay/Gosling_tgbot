package main

import (
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
	yaml "gopkg.in/yaml.v3"
)

func main() {

	botapi, contentPath, listUsers, startMessage, sendMessage := readCfg()[0], readCfg()[1], readCfg()[2], readCfg()[3], readCfg()[4]

	log.Println(botapi, contentPath, listUsers, startMessage, sendMessage)

	TelegramBot(botapi, contentPath, listUsers, startMessage, sendMessage)

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
	usersContent := (cfgYaml["files"].(map[string]interface{})["users_content"])
	usersList := (cfgYaml["files"].(map[string]interface{})["users_list"])

	startMsg := (cfgYaml["messages"].(map[string]interface{})["start_message"])
	sendMsg := (cfgYaml["messages"].(map[string]interface{})["send_message"])

	apiKey_ := fmt.Sprintf("%v", apiKey)
	usersContent_ := fmt.Sprintf("%v", usersContent)
	usersList_ := fmt.Sprintf("%v", usersList)
	startMsg_ := fmt.Sprintf("%v", startMsg)
	sendMsg_ := fmt.Sprintf("%v", sendMsg)

	var out []string
	out = append(out, apiKey_, usersContent_, usersList_, startMsg_, sendMsg_)

	return out
}
