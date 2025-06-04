package quiz

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/neo451/alpha/internal/characters"
	"github.com/neo451/alpha/internal/config"
	"math/rand"
	"os"
	"strings"
	"text/template"
)

type model struct {
	cfg         config.Config
	chars       []characters.Character
	char        characters.Character
	textInput   textinput.Model
	attempts    int
	correct     int
	quitting    bool
	showMessage string
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) RenderOk() string {
	return m.cfg.Theme.CorrectColor.Render(m.cfg.Prompt.Ok)
}

func (m model) RenderErr() string {
	str := fmt.Sprintf(m.cfg.Prompt.Err, m.char.Spelling)
	return m.cfg.Theme.IncorrectColor.Render(str)
}

func (m model) RenderProgress() string {
	progressStr := "\nProgress: 0/0 (0%)\n\n"
	if m.attempts > 0 {
		progressStr = fmt.Sprintf("\nProgress: %d/%d (%.0f%%)\n",
			m.correct, m.attempts, float64(m.correct)/float64(m.attempts)*100)
	}
	return m.cfg.Theme.ProgressColor.Render(progressStr)
}

func (m model) RenderPrompt() string {
	promptStr, err := RenderTemplate(m.cfg.Prompt.Format, m.char)
	if err != nil {
		promptStr = "[template error]\n"
	}

	return m.cfg.Theme.PromptColor.Render(promptStr)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			input := strings.TrimSpace(m.textInput.Value())

			if strings.EqualFold(input, m.cfg.Cmd.Exit) {
				m.quitting = true
				return m, tea.Quit
			}

			m.attempts++
			if strings.EqualFold(input, m.char.Spelling) {
				m.correct++
				m.showMessage = m.RenderOk()
			} else {
				m.showMessage = m.RenderErr()
			}

			m.textInput.SetValue("")
			m.char = m.chars[rand.Intn(len(m.chars))] // Next prompt
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) RenderLeave() string {
	str := fmt.Sprintf("\nFinal score: %d/%d (%.0f%%)\nGoodbye!\n", m.correct, m.attempts, float64(m.correct)/float64(m.attempts)*100)
	return m.cfg.Theme.QuitColor.Render(str)
}

func (m model) RenderWelcome() string {
	str := fmt.Sprintf("Alpha v0.1 — Type '%s' to quit", m.cfg.Cmd.Exit)
	return m.cfg.Theme.WelcomeColor.Render(str)
}

func (m model) View() string {
	if m.quitting {
		return m.RenderLeave()
	}

	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n",
		m.RenderWelcome(),
		m.RenderPrompt(),
		m.textInput.View(),
		m.RenderProgress(),
		m.showMessage)
}

func RenderTemplate(tmplStr string, ctx any) (string, error) {
	tmpl, err := template.New("tpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func initialModel(cfg config.Config, chars []characters.Character) model {
	ti := textinput.New()
	ti.Placeholder = "Type here"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return model{
		cfg:       cfg,
		chars:     chars,
		char:      chars[rand.Intn(len(chars))],
		textInput: ti,
	}
}

// Call this from your main function
func Loop(cfg config.Config, chars []characters.Character) {
	p := tea.NewProgram(initialModel(cfg, chars))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
