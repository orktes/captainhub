package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/google/go-github/github"
)

type captainConfigPlugin struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

type captainConfig struct {
	Plugins []captainConfigPlugin `json:"plugins"`
}

func parseCaptainConf(content []byte, cfg *captainConfig) (err error) {
	return json.Unmarshal(content, cfg)
}

func getCaptainConfig(owner string, repo string) (cfg *captainConfig, err error) {
	cfg = &captainConfig{}

	client := getGithubClient()
	fileContent, _, resp, err := client.Repositories.GetContents(
		owner,
		repo, ".captain.conf", &github.RepositoryContentGetOptions{})

	if err != nil && resp.StatusCode == 404 {
		if err != nil {
			fmt.Printf("Could not load .captain.conf: %s %d\n", err.Error(), resp.StatusCode)
		} else {
			fmt.Printf("Could not load .captain.conf %d\n", resp.StatusCode)
		}
		err = nil
		fileContent = nil
	} else if err != nil {
		return nil, err
	}

	if err == nil && fileContent != nil {
		str, _ := base64.StdEncoding.DecodeString(*fileContent.Content)
		err = parseCaptainConf(str, cfg)
	}

	return
}

func getCaptainPlugin(owner string, repo string, pluginName string) (data []byte, err error) {
	data, err = Asset(pluginName + ".js")

	if err == nil {
		return
	}

	client := getGithubClient()
	fileContent, _, _, err := client.Repositories.GetContents(
		owner,
		repo, "plugins/"+pluginName+".js", &github.RepositoryContentGetOptions{})

	if err != nil {
		return nil, err
	}

	if err == nil && fileContent != nil {
		data = []byte(*fileContent.Content)
	}

	return
}
