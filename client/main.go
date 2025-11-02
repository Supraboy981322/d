package main

import (
	"net/http"
	"log"
	"strings"
	"time"
	"os"
	"io"

	"github.com/BurntSushi/toml"
)

var (
	url string
	line string
	confDir = "/.config/Supraboy981322/d"
	confPath = "config.toml"
	defaultConfig = []byte(`[server]
address = "https://example.com"`)
)

type (
	Config struct {
		Addr string `toml:"address"`
	}
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	if len(os.Args) > 1 {
		line = os.Args[1]
	} else {
		log.Printf("no input")
		return
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(line))
	if err != nil {
		log.Fatalf("err sending request:  %v\n", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("err sending request: %v\n", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read body:  %v\n", err)
	}

	log.Println(string(body))

	defer resp.Body.Close()
}

func readConf(what string) string {
	homeDir, err := os.UserHomeDir()
	hanFerr(err)

	confDir = homeDir + confDir
	confPath = confDir + confPath

	_, err = os.Stat(confPath)
	if os.IsNotExist(err) {
		hanErr(os.MkdirAll(confDir, 0755))
		
		hanFrr(os.WriteFile(confPath, defaultConfig, 0644))
	}

	var conf Config
	_, err := toml.DecodeFile(confPath, &conf)
	hanFrr(err)

	url = conf.Addr
	if url == "https://example.com/" {
		wserr("\033[31m..you don't appear to " + 
						"have configured the address" + 
						"for your server\033[0m")
		fserr("....see \033[32m-h\033[0m")
	}	
}
