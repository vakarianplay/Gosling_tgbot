package main

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func main() {

	botapi, textPath, picDir, userFile, ffont, fsize := readCfg()[0], readCfg()[1], readCfg()[2], readCfg()[3], readCfg()[4], readCfg()[5]
	picEntry(textPath, picDir, ffont, fsize)
	telegramBot(botapi, userFile)

	// fmt.Println(botapi, textPath, picDir, userFile)

}

func readCfg() []string {

	var cfgYaml map[string]interface{}
	cfgFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(cfgFile, &cfgYaml)

	if err != nil {
		log.Fatal(err)
	}

	apiKey := (cfgYaml["bot"].(map[string]interface{})["api_key"])
	textPath := (cfgYaml["files"].(map[string]interface{})["text"])
	picDir := (cfgYaml["files"].(map[string]interface{})["pic_dir"])
	userFile := (cfgYaml["files"].(map[string]interface{})["users_list"])
	font := (cfgYaml["font"].(map[string]interface{})["font_file"])
	fsize := (cfgYaml["font"].(map[string]interface{})["font_size"])

	apiKey_ := fmt.Sprintf("%v", apiKey)
	textPath_ := fmt.Sprintf("%v", textPath)
	picDir_ := fmt.Sprintf("%v", picDir)
	userFile_ := fmt.Sprintf("%v", userFile)
	font_ := fmt.Sprintf("%v", font)
	fsize_ := fmt.Sprintf("%v", fsize)

	var out []string
	out = append(out, apiKey_, textPath_, picDir_, userFile_, font_, fsize_)

	return out
}
