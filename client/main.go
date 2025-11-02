package main

import (
	"net/http"
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
	confPath = "/config.toml"
	defaultConfig = []byte(`[server]
address = "https://example.com/"`)
)

type (
	ServerConfig struct {
		Addr string `toml:"address"`
	}
	Config struct {
		Server ServerConfig
	}
)

func main() {
	readConf()
	if len(os.Args) > 1 {
		line = os.Args[1]
	} else {
		wrl("no input")
		return
	}
	
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(
		"POST", url, strings.NewReader(line))
	hanFrr(err)

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	hanFrr(err)

	body, err := io.ReadAll(resp.Body)
	hanFrr(err)

	wrl(string(body))

	defer resp.Body.Close()
}

func readConf() {
	homeDir, err := os.UserHomeDir()
	hanFrr(err)

	confDir = homeDir + confDir
	confPath = confDir + confPath

	_, err = os.Stat(confPath)
	if os.IsNotExist(err) {
		hanErr(os.MkdirAll(confDir, 0755))
		
		hanFrr(os.WriteFile(
			confPath, defaultConfig, 0644))
	}

	var conf Config
	_, err = toml.DecodeFile(confPath, &conf)
	hanFrr(err)
	
	url = conf.Server.Addr
	if url == "https://example.com/" {
		wserr("\033[31m..you don't appear to " + 
						"have configured the address " + 
						"for your server\033[0m")
		fserr("....see \033[32m-h\033[0m")
	} else if url == "" {
		fserr("shit. something went wrong.\n" +
						"failed to get the address of " +
						"your sever from your config.")
	}
}
