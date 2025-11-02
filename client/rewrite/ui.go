package main

import(
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/bubbles/spinner"
//	"github.com/charmbracelet/lipgloss"
)

/*func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
}*/

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case argsMsg:
		line = string(msg)
		cmd := func() tea.Msg {
			return readConf()()
		}
		return m, cmd
	case confMsg:
		m.url = string(msg)
		cmd := func() tea.Msg {
			return sendReq()()
		}
		return m, cmd
	case respMsg:
		m.response = string(msg)
		return m, tea.Quit
	case errMsg:
		m.err = error(msg)
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.quitting {
		return "quitting...\n"
	} else {
		out := ""
		out += fmt.Sprintf("resp:  %s\n", m.response)
		
		if m.config != nil {
			out += fmt.Sprintf("server:  %s\n", m.config.Server)
		}
	
		if m.err != nil {
			out += fmt.Sprintf("err:  %v\n", m.err)
			out += "\npress q to quit.\n"
			return out
		}
	
		if m.line != "" {
			out += "press q to quit.\n"
			return out
		}
		out += "working...\n"
		return out
	}
}
