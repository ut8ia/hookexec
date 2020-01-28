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
	Request struct {
		Header string `yaml:"header"`
		Token  string `yaml:"token"`
		Param  string `yaml:"param"`
	} `yaml:"request"`
	Hooks map[string]struct {
		Executor   string `yaml:"executor"`
		ScriptPath string `yaml:"scriptPath"`
		Script     string `yaml:"script"`
	} `yaml:"hooks"`
}

func main() {
	configFile := "./configs/config.yml"
	if len(os.Args) >1 {
		configFile = os.Args[1]
	}
	log.Println(configFile)
	readConfig(&cfg, configFile)
	http.HandleFunc("/", RequestHandler)
	serve := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := http.ListenAndServe(serve, nil); err != nil {
		log.Fatal(err)
	}
}

func readConfig(cfg *Config, configFile string) {
	f, err := os.Open(configFile)
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

func RequestHandler(w http.ResponseWriter, r *http.Request) {

	header := r.Header.Get(cfg.Request.Header)
	if header != cfg.Request.Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, exist := r.URL.Query()[cfg.Request.Param]
	var hook = "default"
	if exist {
		_, ok := cfg.Hooks[params[0]]
		if ok {
			hook = params[0]
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, cfg.Server.BodyLimit)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	cmd := exec.Command(cfg.Hooks[hook].Executor, cfg.Hooks[hook].Script, bodyString)
	out, err := cmd.Output()

	if err == nil {
		log.Println(string(out))
	} else {
		log.Println(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
