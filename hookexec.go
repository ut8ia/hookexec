package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var cfg Config

type Config struct {
	Server struct {
		Port      string `yaml:"port"`
		Host      string `yaml:"host"`
		BodyLimit int64  `yaml:"bodyLimit"`
	} `yaml:"server"`
	Auth struct {
		Header    string `yaml:"header"`
		Token     string `yaml:"token"`
	} `yaml:"auth"`
	Default struct {
		Executor   string `yaml:"executor"`
		ScriptPath string `yaml:"scriptPath"`
		Script    string `yaml:"script"`
	} `yaml:"Default"`
}

func main() {
	readConfig(&cfg)
	http.HandleFunc("/", RequestHandler)

	serve := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := http.ListenAndServe(serve, nil); err != nil {
		log.Fatal(err)
	}
}

func readConfig(cfg *Config) {
	f, err := os.Open("./configs/config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {

	var script string

	// parse header and check token
	header := r.Header.Get(cfg.Auth.Header)
	if header != cfg.Auth.Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Double check it's a post request being made
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	commands, ok := r.URL.Query()["hook"]

	if !ok {
		script = cfg.Default.Script
	} else {
		script = commands[0]
	}

	r.Body = http.MaxBytesReader(w, r.Body, cfg.Server.BodyLimit)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	cmd := exec.Command(cfg.Default.Executor, script, bodyString)
	cmd.Dir = cfg.Default.ScriptPath
	out, err := cmd.Output()

	// Log all data. Form is a map[]
	if err == nil {
		log.Println(string(out))
	} else {
		log.Println(err)
	}
}
