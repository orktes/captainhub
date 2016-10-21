package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/google/go-github/github"
	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/labstack/echo"
	mw "github.com/orktes/captainhub/Godeps/_workspace/src/github.com/labstack/echo/middleware"
	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/robertkrimen/otto"
	"github.com/orktes/captainhub/Godeps/_workspace/src/github.com/ryanuber/go-glob"
	"github.com/orktes/captainhub/Godeps/_workspace/src/golang.org/x/oauth2"

	_ "github.com/orktes/captainhub/Godeps/_workspace/src/github.com/robertkrimen/otto/underscore"
)

//go:generate go-bindata -prefix=plugins/ -pkg=main plugins/...

var redisPool *redis.Pool
var hookSecret string
var errMissingSig = echo.NewHTTPError(http.StatusForbidden, "Missing X-Hub-Signature")
var errInvalidSig = echo.NewHTTPError(http.StatusForbidden, "Invalid X-Hub-Signature")

func matchFilePath(call otto.FunctionCall) otto.Value {
	pattern := call.Argument(0).String()
	name := call.Argument(1).String()

	match := glob.Glob(pattern, name)
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

	var owner string
	var repo string

	r := c.Request()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if hookSecret != "" {
		sig := r.Header.Get("X-Hub-Signature")
		if sig == "" {
			return errMissingSig
		}

		mac := hmac.New(sha1.New, []byte(hookSecret))
		mac.Write(body)
		expectedMAC := mac.Sum(nil)
		expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)
		if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
			return errInvalidSig
		}
	}

	eventType := r.Header.Get("X-Github-Event")

	data := map[string]interface{}{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("%s", err.Error())
		return err
	}

	repository := data["repository"]

	repoMap, ok := repository.(map[string]interface{})
	if !ok {
		fmt.Printf("Repository missing")
		return nil
	}

	repo = repoMap["name"].(string)

	ownerMap, ok := repoMap["owner"].(map[string]interface{})
	if !ok {
		fmt.Printf("Owner missing")
		return nil
	}

	if _, ok := ownerMap["login"]; ok {
		owner = ownerMap["login"].(string)
	} else if _, ok := ownerMap["name"]; ok {
		owner = ownerMap["name"].(string)
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
		vm.Run(`
			var global = {};
			_moduleCache = {};
			function require(moduleName) {
				if (!_moduleCache[moduleName]) {
					_require(moduleName);
				}

				return _moduleCache[moduleName];
			}
		`)
		vm.Set("eventType", eventType)
		vm.Set("eventData", data)
		vm.Set("config", plugin.Config)
		vm.Set("matchFilePath", matchFilePath)
		vm.Set("_require", func(call otto.FunctionCall) otto.Value {
			moduleFilename, err := call.Argument(0).ToString()
			if err != nil {
				panic(err)
			}

			content, err := getCaptainPlugin(owner, repo, moduleFilename)
			if err != nil {
				panic(err)
			}
			script := `
				(function (_moduleCache, global) {
				  var module = {exports: {}};
					(function (require, module, exports, undefined) {
						` + string(content) + `
					})(require, module, module.exports);
					_moduleCache["` + moduleFilename + `"] = module.exports;
				})(_moduleCache, global);
			`
			_, err = vm.Run(script)
			if err != nil {
				fmt.Printf("Error evaluating %s %s %s\n", moduleFilename, script, err.Error())
			}

			return otto.Value{}
		})

		vm.Set("print", func(call otto.FunctionCall) otto.Value {
			text, err := call.Argument(0).ToString()
			if err != nil {
				panic(err)
			}

			println(text)
			return otto.Value{}
		})

		vm.Set("getPullRequestFileContent", func(call otto.FunctionCall) otto.Value {
			prNumber, err := call.Argument(0).ToInteger()
			if err != nil {
				panic(err)
			}

			fileName, err := call.Argument(1).ToString()
			if err != nil {
				panic(err)
			}

			content, err := readPullRequestFileContent(owner, repo, int(prNumber), fileName)
			if err != nil {
				panic(err)
			}

			ottoValue, err := otto.ToValue(string(content))
			if err != nil {
				panic(err)
			}

			return ottoValue
		})

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
			redisClient := redisPool.Get()
			defer redisClient.Close()

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
			redisClient := redisPool.Get()
			defer redisClient.Close()

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

		_, err = vm.Run(`
			var module = {exports: {}};
				(function (require, module, exports, global, undefined) {
					` + string(pluginData) + `
				})(require, module, module.exports, global);
			`)
		if err != nil {
			fmt.Printf("plugin eval error: %s %s\n", err.Error(), pluginData)
			return err
		}
	}

	return nil
}

func newRedisPool(serverURL, db string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(serverURL)
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			diff := time.Now().Sub(t)
			// if the client has not been used for 60 seconds ping it before handing it out
			if diff.Seconds() >= 60 {
				_, err := c.Do("PING")
				return err
			}
			return nil
		},
		Wait: true,
	}
}

func main() {
	// Init hook secret
	hookSecret = os.Getenv("GITHUB_HOOK_SECRET")

	// Init redis client
	redisPool = newRedisPool(os.Getenv("REDISCLOUD_URL"), "0")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	//e.Use(mw.Recover())

	// Routes
	e.Post("/payload", payload)

	// Start server
	e.Run(":" + os.Getenv("PORT"))
}
