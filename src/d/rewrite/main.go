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
		spinOn bool
		quitting bool
		chkErr bool
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

func chkArgs(m model) tea.Cmd {
	m.spinOn = true
	return func() tea.Msg {
		//chk args
		if len(os.Args) > 1 {
			line = os.Args[1]
			return argsMsg(line)
		}
		return quitErrMsg(merr("no input", nil))
	}
}

func readConf(m model) tea.Cmd {
	m.spinOn = true
	return func() tea.Msg {
		//get ~
		homeDir, err := os.UserHomeDir()
		if err != nil {	
			return quitErrMsg(merr("failed to get home dir", nil))
		}

		//set actual conf dir and path
		confDir = homeDir + confDir
		confPath = confDir + confPath
	
		//chk if it exists
		_, err = os.Stat(confPath)
		if os.IsNotExist(err) {
			//mk if not
			err := os.MkdirAll(confDir, 0755)
			if err != nil {
				return quitErrMsg(merr("failed to create config dir", nil))
			}
			
			//write default conf
			err = (os.WriteFile(
				confPath, defaultConfig, 0644))
			if err != nil {
				return quitErrMsg(merr("failed to create config", nil))
			}
		}
	
		//read conf
		var conf Config
		_, err = toml.DecodeFile(confPath, &conf)
		if err != nil {
			return quitErrMsg(merr("failed to decode config", nil))
		}
		
		//set url from conf
		url := conf.Server.Addr
		if url == "https://example.com/" || url == "" {
			return quitErrMsg(merr("server address not set", nil))
		}
		return confMsg(url)
	}
}

func sendReq(m model) tea.Cmd {
	return func() tea.Msg {
		//create http client
		//  timeout after 10 seconds
		client := &http.Client{
			Timeout: time.Second * 10,
		}
	
		//create req
			req, err := http.NewRequest(
			"POST", m.url, strings.NewReader(line))
		if err != nil {
			return quitErrMsg(merr("failed to create request\n", nil))
		}
	
		//set header (not needed, but just to be safe)
		req.Header.Set("Content-Type", "text/plain")

		//mk req
		resp, err := client.Do(req)
		if err != nil {
			wrl(url)
			return quitErrMsg(merr("failed to send request\n", nil))
		}
	
		//read resp 
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return quitErrMsg(merr("failed to read response\n", nil))
		}
	
		//return resp
		return respMsg(string(body))
	}
}
