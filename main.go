package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	count    int
	songName string
}

type mpvFinishedMsg struct {
	err error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "a":
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Printf("failed to access user home dir")
			}
			filePath := filepath.Join(homeDir, "/Music/loyalty_card")
			c := exec.Command("mpv", "--vo=kitty", "--vo-kitty-use-shm=yes", filePath, ".")
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				if err != nil {
					log.Fatalf("failed to start")
				}
				return nil
			})
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	return tea.NewView(fmt.Sprintf("The Player im trying to make"))

}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
