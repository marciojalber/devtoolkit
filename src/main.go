package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	quiting bool
	cursor  int
	choices []string

	busy   bool
	timer  time.Time
	report string
}

func initModel() model {
	m := model{
		choices: []string{
			"To compile",
			"To run",
			"To doc",
		},
	}

	return m
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(data tea.Msg) (tea.Model, tea.Cmd) {
	if !m.busy {
		m.report = ""
	}

	switch data := data.(type) {
	case time.Time:
		msg := "Processing. Please wait... (%d)"
		secs := int(time.Since(m.timer).Seconds())
		m.report = fmt.Sprintf(msg, secs)

		if secs < 5 {
			return m, tick()
		}

		m.busy = false
		m.report = ""

	case tea.KeyMsg:
		switch data.String() {
		case "esc", "ctrl+c":
			m.quiting = true
			return m, tea.Quit
		}

		if m.busy {
			break
		}

		switch data.String() {
		case "up", "w":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "s":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.busy = true
			m.timer = time.Now()
			m.report = "Processing. Please wait..."
			return m, tick()
		}
	}

	return m, nil
}

func (m model) View() string {
	content := []string{
		"",
		"=============================",
		"",
		"     Devtoolkit started.",
		"",
		"=============================",
		":: Press \"Esc\" or \"Ctrl+c\" to quit",
		"",
		"Which you wanna do?",
		"",
	}

	for i, choice := range m.choices {
		var cursor = "  "
		if m.cursor == i {
			cursor = "->"
		}
		content = append(content, cursor+" "+choice)
	}
	content = append(content, "")

	if m.quiting {
		content = append(content, "Bye!", "")
	} else {
		content = append(content, m.report)
	}

	return strings.Join(content, "\n")
}

func main() {
	program := tea.NewProgram(initModel())
	if _, err := program.Run(); err != nil {
		fmt.Println("Error: " + err.Error())
	}
}
