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
	url string   //set later when read conf
	line string
	timeout int

	//changes to actual path
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
		Timeout int `toml:"timeout"`
		Server ServerConfig
	}
)

func main() {
	//fn to read conf (clearly)
	readConf()
	
	//chk args
	if len(os.Args) > 1 {
		line = strings.Join(os.Args[1:], " ")
	} else {
		wrl("no input")
		return
	}
	
	//create http client
	//  timeout after 10 seconds
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	//create req
	req, err := http.NewRequest(
		"POST", url, strings.NewReader(line))
	hanFrr(err)

	//set header (not needed, but just to be safe)
	req.Header.Set("Content-Type", "text/plain")

	//mk req
	resp, err := client.Do(req)
	hanFrr(err)

	//read resp 
	body, err := io.ReadAll(resp.Body)
	hanFrr(err)

	//print resp
	wrl(string(body))

	//why do I have to call this fn?
	//  it's literally the last thing done
	//    why can't it just close itself in
	//      scenario
	defer resp.Body.Close()
}

func readConf() {
	//get ~
	homeDir, err := os.UserHomeDir()
	hanFrr(err)

	//set actual conf dir and path
	confDir = homeDir + confDir
	confPath = confDir + confPath

	//chk if it exists
	_, err = os.Stat(confPath)
	if os.IsNotExist(err) {
		//mk if not
		hanErr(os.MkdirAll(confDir, 0755))
		
		//write default conf
		hanFrr(os.WriteFile(
			confPath, defaultConfig, 0644))
	}

	//read conf
	var conf Config
	_, err = toml.DecodeFile(confPath, &conf)
	hanFrr(err)
	
	//set url from conf
	url = conf.Server.Addr

	timeout = conf.Timeout

	//blame user for all other problems
	if url == "https://example.com/" {
		wserr("\033[31m..you don't appear to " + 
						"have configured the address " + 
						"for your server\033[0m")
		fserr("....see \033[32m-h\033[0m")
	} else if url == "" {
		//unless if it's blank, then it's assumed
		//  to be a problem with the client code
		fserr("shit. something went wrong.\n" +
						"failed to get the address of " +
						"your sever from your config.")
	}
}
