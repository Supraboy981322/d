package main

import (
//	"fmt"
	"net/http"
	"strings"
	"time"
	"os"
	"io"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/BurntSushi/toml"
)

var (
	url string   //set later when read conf
	line string

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
		Server ServerConfig
	}

	confMsg string
	argsMsg string
	respMsg string
	errMsg struct { error }
	quitErrMsg error

	model struct {
		url string
		config *Config
		line string
		response string
		spinner spinner.Model
		quitting bool
		noInput bool
		err error
	}
)

func (m model) Init() tea.Cmd {
	return chkArgs()
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		ferr(err)
	}
}

func (e errMsg) Error() string {
	return e.error.Error()
}

func chkArgs() tea.Cmd {
	return func() tea.Msg {
		//chk args
		if len(os.Args) > 1 {
			line = os.Args[1]
			return argsMsg(line)
		}
		return quitErrMsg(merr("no input", nil))
	}
}

func readConf() tea.Cmd {
	return func() tea.Msg {
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
		url := conf.Server.Addr
		if url == "https://example.com/" || url == "" {
			if e, ok := merr("server address not set", nil).(errMsg); ok {
				return errMsg(e)
			}
		}
		return confMsg(url)
	}
}

func sendReq() tea.Cmd {
	return func() tea.Msg {
		//create http client
		//  timeout after 10 seconds
		client := &http.Client{
			Timeout: time.Second * 10,
		}
	
		//create req
			req, err := http.NewRequest(
			"POST", url, strings.NewReader(line))
		if err != nil {
			if e, ok := merr("server address not set", err).(errMsg); ok {
				return errMsg(e)
			}
		}
	
		//set header (not needed, but just to be safe)
		req.Header.Set("Content-Type", "text/plain")
	
		//mk req
		resp, err := client.Do(req)
		if err != nil {
			if e, ok := merr("failed to send request", err).(errMsg); ok {
				return errMsg(e)
			}
		}
	
		//read resp 
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if e, ok := merr("server address not set", err).(errMsg); ok {
				return errMsg(e)
			}
		}
	
		//return resp
		return respMsg(string(body))
	}
}
