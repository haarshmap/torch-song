package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ioksc/gompv"
)

type model struct {
	count int
	mpv   *gompv.Client
}

type mpvFinishedMsg struct {
	err error
}

var (
	container = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#cc6666")).
		Width(120).
		Height(39)
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j":
			m.mpv.Next()
			return m, nil

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	v := tea.NewView(fmt.Sprintf(container.Render("The Player im trying to make")))
	v.AltScreen = true

	return v
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to access user home dir")
	}
	filePath := filepath.Join(homeDir, "/Music/loyalty_card")
	socketPath := "/tmp/mpv.sock"
	server, client, err := gompv.StartAndConnect(context.Background(), socketPath, filePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Socket connection established")
	defer server.Stop()
	defer client.Close()

	p := tea.NewProgram(model{mpv: client})
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}

}
