package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate datasync.yaml",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
		// Bubble teaを使って、UIを作る。
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type ScreenType int

const (
	SelectTargetKind ScreenType = iota
	InputMysql
	InputFile
	ConfirmAddTarget
	ConfirmSetupRemote
	SelectRemoteKind
	InputGcs
)

type model struct {
	screenType ScreenType

	// 入力共有
	focusIndex int
	inputs     []textinput.Model

	// TODO
}

func initialModel() model {
	m := model{
		screenType: SelectTargetKind,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screenType {
	case SelectTargetKind:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "up":
				if m.focusIndex == 1 {
					m.focusIndex = 0
				}
			case "down":
				if m.focusIndex == 0 {
					m.focusIndex = 1
				}
			case "enter":
				m.screenType = InputMysql
				m.focusIndex = 0
				m.inputs = makeMysqlInputs()
			}
		}
	case InputMysql:
		cmds := make([]tea.Cmd, len(m.inputs))
		for i := 0; i <= len(m.inputs)-1; i++ {
			if i == m.focusIndex {
				// Set focused state
				cmds[i] = m.inputs[i].Focus()
				m.inputs[i].PromptStyle = focusedStyle
				m.inputs[i].TextStyle = focusedStyle
				continue
			}
			// Remove focused state
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
		return m, tea.Batch(cmds...)
	default:
		panic("invalid screenType")

	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func makeMysqlInputs() []textinput.Model {
	var inputs []textinput.Model
	var t textinput.Model
	inputs = make([]textinput.Model, 0)

	// host:
	// 	localhost
	// port:
	// 	3300
	// user:
	// 	root
	// password:
	// 	root
	// database:
	// 	sample

	t = textinput.New()
	t.Cursor.Style = cursorStyle
	t.Placeholder = "Hostname"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	inputs = append(inputs, t)

	t = textinput.New()
	t.Cursor.Style = cursorStyle
	t.Placeholder = "port"
	t.PromptStyle = blurredStyle
	t.TextStyle = blurredStyle
	inputs = append(inputs, t)

	return inputs
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	switch m.screenType {
	case SelectTargetKind:
		b.WriteString("? How kind of dump target? …\n")
		if m.focusIndex == 0 {
			b.WriteString("❯ MySQL\n")
			b.WriteString("  File(s)\n")
		} else {
			b.WriteString("  MySQL\n")
			b.WriteString("❯ File(s)\n")
		}
	case InputMysql:
		b.WriteString("Input mysql setting …\n")
		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				b.WriteRune('\n')
			}
		}
	}

	return b.String()
}

/*
? How kind of dump target? …
❯ MySQL
  File(s)

-- mysql
? MySQL server hostname / port / username / password / databasename
>

-- file
? Select directory or file
> picker


? Add dump target?
  Yes
❯ No


? Setup remote server?
❯ Yes
  No

? Remote server type?
❯ Google Cloud Storage
  Amazon S3
	Samba

-- GCS
? GCS bucket / path
>

*/
