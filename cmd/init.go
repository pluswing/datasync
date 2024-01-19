package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pluswing/datasync/data"
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

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

type model struct {
	screenType ScreenType

	// 入力共有
	focusIndex int
	inputs     []textinput.Model

	// ファイル選択
	filepicker filepicker.Model
	err        error

	targets []interface{}
}

func initialModel() model {
	m := model{
		screenType: SelectTargetKind,
		targets:    make([]interface{}, 0),
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screenType {
	case SelectTargetKind:
		return updateSelectTargetKind(m, msg)
	case InputMysql:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "enter", "up", "down":
				s := msg.String()
				if s == "up" {
					m.focusIndex -= 1
					if m.focusIndex < 0 {
						m.focusIndex = 0
					}
				} else if s == "down" {
					m.focusIndex += 1
					if m.focusIndex >= len(m.inputs) {
						m.focusIndex = len(m.inputs) - 1
					}
				} else if s == "enter" {
					if m.focusIndex == len(m.inputs)-1 {
						port, err := strconv.Atoi(m.inputs[1].Value())
						if err != nil {
							fmt.Println("invalut port")
						}
						var t = data.TargetType{
							Kind: "mysql",
							Config: data.TargetMysqlType{
								Host:     m.inputs[0].Value(),
								Port:     port,
								User:     m.inputs[2].Value(),
								Password: m.inputs[3].Value(),
								Database: m.inputs[4].Value(),
							},
						}
						m.targets = append(m.targets, t)
						// 次のスクリーンに行く。
						m.screenType = ConfirmAddTarget
						m.focusIndex = 1
						m.inputs = make([]textinput.Model, 0)
					} else {
						m.focusIndex += 1
						if m.focusIndex >= len(m.inputs) {
							m.focusIndex = len(m.inputs) - 1
						}
					}
				}
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
			}
		}
	case InputFile:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		case clearErrorMsg:
			m.err = nil
		}

		var cmd tea.Cmd
		m.filepicker, cmd = m.filepicker.Update(msg)

		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			cwd, _ := os.Getwd()
			path = strings.Replace(path, cwd, "", 1)
			fmt.Printf("path = %s\n", path)
			var t = data.TargetType{
				Kind: "file",
				Config: data.TargetFileType{
					Path: path,
				},
			}
			m.targets = append(m.targets, t)
			m.screenType = ConfirmAddTarget
			m.focusIndex = 1
			m.inputs = make([]textinput.Model, 0)
		}

		if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
			m.err = errors.New(path + " is not valid.")
			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}

		return m, cmd
	case ConfirmAddTarget:
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
				if m.focusIndex == 0 {
					m.screenType = SelectTargetKind
					m.focusIndex = 0
				} else {
					// TODO
					m.screenType = ConfirmSetupRemote
					m.focusIndex = 0
				}
			}
		}
	default:
		panic("invalid screenType")
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func updateSelectTargetKind(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.focusIndex == 0 {
				m.screenType = InputMysql
				m.focusIndex = 0
				m.inputs = makeMysqlInputs()
			} else {
				m.screenType = InputFile
				m.focusIndex = 0
				m.inputs = make([]textinput.Model, 0)
				fp := filepicker.New()
				fp.DirAllowed = true
				fp.CurrentDirectory, _ = os.Getwd()
				fp.Height = 10
				cmd := fp.Init()
				m.filepicker = fp
				return m, cmd
			}
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func makeMysqlInputs() []textinput.Model {
	var inputs []textinput.Model
	var t textinput.Model
	inputs = make([]textinput.Model, 0)

	t = textinput.New()
	t.Cursor.Style = cursorStyle
	t.Placeholder = "Hostname (default: localhost)"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	inputs = append(inputs, t)

	t = textinput.New()
	t.Placeholder = "port (default: 3306)"
	inputs = append(inputs, t)

	t = textinput.New()
	t.Placeholder = "user (default: root)"
	inputs = append(inputs, t)

	t = textinput.New()
	t.Placeholder = "password (default: '')"
	inputs = append(inputs, t)

	t = textinput.New()
	t.Placeholder = "database"
	inputs = append(inputs, t)

	return inputs
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
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
		ViewSelect(&b, m.focusIndex, []string{"MySQL", "File(s)"})
	case InputMysql:
		b.WriteString("Input mysql setting …\n") // FIXME これだけ残る。なんとかする。
		ViewInputs(&b, m.inputs)
	case InputFile:
		if m.err != nil {
			b.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
		} else {
			b.WriteString("Pick a file or directory:\n")
			// b.WriteString("<- : Parent Directory, -> : Dig Directory, enter : select")
		}
		b.WriteString("\n\n" + m.filepicker.View() + "\n")
	case ConfirmAddTarget:
		b.WriteString("? Add dump target? …\n")
		ViewSelect(&b, m.focusIndex, []string{"Yes", "No"})
	}
	return b.String()
}

func ViewSelect(b *strings.Builder, focusIndex int, texts []string) {
	for i, t := range texts {
		if focusIndex == i {
			b.WriteString(focusedStyle.Render(fmt.Sprintf("❯ %s\n", t)))
		} else {
			b.WriteString(fmt.Sprintf("\r  %s\n", t))
		}
	}
}

func ViewInputs(b *strings.Builder, inputs []textinput.Model) {
	for i := range inputs {
		b.WriteString(inputs[i].View())
		if i < len(inputs)-1 {
			b.WriteRune('\n')
		}
	}
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
