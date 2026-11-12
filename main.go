package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ioksc/gompv"
)

type model struct {
	mpv     *gompv.Client
	title   string
	timePos string
	paused  string
}

var (
	container = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#cc6666")).
		Padding(1, 2)
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gompv.Event:
		if name, value, ok := msg.PropertyString(); ok && name == "media-title" {
			m.title = value
		}
		if name, value, ok := msg.PropertyFloat(); ok && name == "time-pos" {
			m.timePos = strconv.FormatFloat(value, 'f', 0, 64)
		}
		if name, value, ok := msg.PropertyBool(); ok && name == "pause" {
			if value == true {
				m.paused = "󰏤"
			}
			if value == false {
				m.paused = ""
			}
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "a":
			m.mpv.TogglePause()
			return m, nil
		case "k":
			m.mpv.Next()
			return m, nil
		case "j":
			m.mpv.Prev()
			return m, nil
		case "esc":
			m.mpv.ToggleMute()
			return m, nil
		case "v":
			m.mpv.Shuffle()
			return m, nil
		case "l":
			total, _ := m.mpv.GetVolume()
			total += 5
			if total > 100 {
				total = 100
			}
			m.mpv.SetVolume(total)
			return m, nil
		case "h":
			total, _ := m.mpv.GetVolume()
			total -= 5
			if total < 0 {
				total = 0
			}
			m.mpv.SetVolume(total)
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	vol, _ := m.mpv.GetVolume()
	volStr := strconv.FormatFloat(vol, 'f', 0, 64)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"torch-song",
		"current song: "+m.title+"\tVolume: "+volStr,
		"",
		m.timePos, m.paused,
		"",
		"j -> prev | k -> next | h/l -> adjust volume | a -> pause/play | capsLock -> mute | v -> shuffle",
	)

	v := tea.NewView(container.Render(content))
	v.AltScreen = true

	return v
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to access user home dir")
	}
	filePath := filepath.Join(homeDir, "/Music")
	socketPath := "/tmp/mpv.sock"
	server, client, err := gompv.StartAndConnect(context.Background(), socketPath, filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Stop()
	defer client.Close()

	err = client.ObserveProperties("media-title", "time-pos", "pause")
	if err != nil {
		log.Printf("Failed to observe properties")
	}

	p := tea.NewProgram(model{mpv: client})
	go func() {
		for ev := range client.Events() {
			p.Send(ev)
		}
	}()

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}

}
