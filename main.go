package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/google/go-github/github"
	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/labstack/echo"
	mw "github.com/orktes/captainhub/Godeps/_workspace/src/github.com/labstack/echo/middleware"
	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/robertkrimen/otto"
	"github.com/orktes/captainhub/Godeps/_workspace/src/golang.org/x/oauth2"

	_ "github.com/orktes/captainhub/Godeps/_workspace/src/github.com/robertkrimen/otto/underscore"
)

//go:generate go-bindata -prefix=plugins/ -pkg=main plugins/...

var redisClient redis.Conn

func matchFilePath(call otto.FunctionCall) otto.Value {
	pattern := call.Argument(0).String()
	name := call.Argument(1).String()

	match, _ := filepath.Match(pattern, name)
	val, _ := otto.ToValue(match)
	return val
}

func getGithubClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	return client
}

// Handler
func payload(c *echo.Context) error {
	r := c.Request()
	eventType := r.Header.Get("X-Github-Event")

	owner := c.Param("owner")
	repo := c.Param("repo")

	data := map[string]interface{}{}
	if err := c.Bind(&data); err != nil {
		fmt.Printf("%s", err.Error())
		return err
	}

	fmt.Printf("Event: %s, Owner: %s, Repo %s\n", eventType, owner, repo)

	cfg, err := getCaptainConfig(owner, repo)

	fmt.Printf("%s/%s config: %+v\n", owner, repo, cfg)

	if err != nil {
		return err
	}

	for _, plugin := range cfg.Plugins {
		pluginData, err := getCaptainPlugin(owner, repo, plugin.Name)

		if err != nil {
			// TODO set info about this
			fmt.Printf("%s", err.Error())
			continue
		}

		vm := otto.New()
		vm.Set("eventType", eventType)
		vm.Set("eventData", data)
		vm.Set("config", plugin.Config)
		vm.Set("matchFilePath", matchFilePath)

		vm.Set("getPullRequestDetails", func(call otto.FunctionCall) otto.Value {
			prNumber, err := call.Argument(0).ToInteger()
			if err != nil {
				panic(err)
			}

			pr, err := getPullRequestDetails(owner, repo, int(prNumber))

			data, err := json.Marshal(pr)
			if err != nil {
				panic(err)
			}

			val, err := vm.Run(fmt.Sprintf("(%s)", string(data)))
			if err != nil {
				panic(err)
			}

			return val
		})

		vm.Set("getPullRequestFiles", func(call otto.FunctionCall) otto.Value {
			prNumber, err := call.Argument(0).ToInteger()
			if err != nil {
				panic(err)
			}

			files, err := listPullRequestFiles(owner, repo, int(prNumber))
			if err != nil {
				panic(err)
			}

			data, err := json.Marshal(files)
			if err != nil {
				panic(err)
			}

			val, err := vm.Run(string(data))
			if err != nil {
				panic(err)
			}

			return val
		})

		vm.Set("createPullRequestComment", func(call otto.FunctionCall) otto.Value {
			prNumber, err := call.Argument(0).ToInteger()
			if err != nil {
				panic(err)
			}

			body := call.Argument(1).String()
			if err != nil {
				panic(err)
			}

			err = createPullRequestComment(owner, repo, int(prNumber), body)
			if err != nil {
				panic(err)
			}

			return otto.Value{}
		})

		vm.Set("createIssueComment", func(call otto.FunctionCall) otto.Value {
			prNumber, err := call.Argument(0).ToInteger()
			if err != nil {
				panic(err)
			}

			body := call.Argument(1).String()
			if err != nil {
				panic(err)
			}

			err = createIssueComment(owner, repo, int(prNumber), body)
			if err != nil {
				panic(err)
			}

			return otto.Value{}
		})

		vm.Set("createStatus", func(call otto.FunctionCall) otto.Value {
			sha, err := call.Argument(0).ToString()
			if err != nil {
				panic(err)
			}

			state, err := call.Argument(1).ToString()
			if err != nil {
				panic(err)
			}

			targetURL, err := call.Argument(2).ToString()
			if err != nil {
				panic(err)
			}

			description, err := call.Argument(3).ToString()
			if err != nil {
				panic(err)
			}

			context, err := call.Argument(4).ToString()
			if err != nil {
				panic(err)
			}

			err = createStatus(owner, repo, sha, state, targetURL, description, context)
			if err != nil {
				panic(err)
			}

			return otto.Value{}
		})

		vm.Set("saveData", func(call otto.FunctionCall) otto.Value {
			key := call.Argument(0).String()
			data := call.Argument(1).String()

			storeKey := fmt.Sprintf("%s/%s/%s", owner, repo, key)

			_, err := redisClient.Do("SET", storeKey, data)
			if err != nil {
				panic(err)
			}

			return otto.Value{}
		})

		vm.Set("loadData", func(call otto.FunctionCall) otto.Value {
			key := call.Argument(0).String()
			storeKey := fmt.Sprintf("%s/%s/%s", owner, repo, key)

			s, err := redis.String(redisClient.Do("GET", storeKey))
			if err != nil {
				panic(err)
			}

			ottoValue, err := otto.ToValue(s)
			if err != nil {
				panic(err)
			}

			return ottoValue
		})

		_, err = vm.Run(string(pluginData))
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return err
		}
	}

	return nil
}

func main() {
	// Init redis client
	var err error
	redisClient, err = redis.DialURL(os.Getenv("REDISCLOUD_URL"))
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Post("/:owner/:repo", payload)

	// Start server
	e.Run(":" + os.Getenv("PORT"))
}
