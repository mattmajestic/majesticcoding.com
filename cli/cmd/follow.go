package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var followCmd = &cobra.Command{
	Use:   "follow",
	Short: "Select socials to follow",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tea.NewProgram(initialModel()).Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(followCmd)
}

// Bubble Tea Model
type model struct {
	cursor   int
	selected map[int]bool
	options  []string
	message  string
}

// Init method for the Bubble Tea model
func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	return model{
		cursor:   0,
		selected: make(map[int]bool),
		options:  []string{"GitHub", "YouTube", "Twitch"},
		message:  "Use ↑/↓ to navigate, Space/Enter to select, 'o' to open, 'q' to quit.",
	}
}

// Bubble Tea Update Function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q": // Quit
			return m, tea.Quit
		case "up", "k": // Move cursor up
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j": // Move cursor down
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter", " ": // Select or deselect
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "o": // Open selected socials
			for i, selected := range m.selected {
				if selected {
					openSocial(m.options[i])
				}
			}
			m.message = "Opened selected socials in your browser!"
		}
	}
	return m, nil
}

// Bubble Tea View Function
func (m model) View() string {
	var ui string
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	ui += titleStyle.Render("Select socials to follow:\n\n")
	for i, option := range m.options {
		cursor := " " // No cursor by default
		if m.cursor == i {
			cursor = cursorStyle.Render(">")
		}

		checkbox := "[ ]" // Unchecked by default
		if m.selected[i] {
			checkbox = selectedStyle.Render("[x]")
		}

		ui += fmt.Sprintf("%s %s %s\n", cursor, checkbox, option)
	}
	ui += "\n" + messageStyle.Render(m.message)
	return ui
}

// Open Social Links
func openSocial(name string) {
	socials := map[string]string{
		"GitHub":  "https://github.com/mattmajestic",
		"YouTube": "https://www.youtube.com/channel/UCjTavL86-CW6j58fsVIjTig?sub_confirmation=1",
		"Twitch":  "https://www.twitch.tv/MajesticCodingTwitch",
	}

	url, exists := socials[name]
	if !exists {
		fmt.Printf("No link available for %s\n", name)
		return
	}

	fmt.Printf("Opening %s...\n", name)
	_ = browser.OpenURL(url)
}
