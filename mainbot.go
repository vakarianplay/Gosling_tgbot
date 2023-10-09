package main

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func main() {

	fmt.Println(readCfg()[0])
	picEntry(readCfg()[1], readCfg()[2])
}

func readCfg() []string {

	var cfgYaml map[string]interface{}
	cfgFile, err := os.ReadFile("config.yml")
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

	apiKey_ := fmt.Sprintf("%v", apiKey)
	textPath_ := fmt.Sprintf("%v", textPath)
	picDir_ := fmt.Sprintf("%v", picDir)

	var out []string
	out = append(out, apiKey_, textPath_, picDir_)

	return out
}
