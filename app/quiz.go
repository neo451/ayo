package app

// quiz: match user spelling to answer

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/template"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/neo451/ayo/internal/characters"
	"github.com/neo451/ayo/internal/config"
)

type quiz struct {
	cfg         config.Config
	chars       []characters.Character
	char        characters.Character
	textInput   textinput.Model
	attempts    int
	correct     int
	quitting    bool
	showMessage string
}

func (m quiz) Init() tea.Cmd {
	return textinput.Blink
}

func (m quiz) RenderOk() string {
	return m.cfg.Theme.CorrectColor.Render(m.cfg.Prompt.Ok)
}

func (m quiz) RenderErr() string {
	str := fmt.Sprintf(m.cfg.Prompt.Err, m.char.Symbol, m.char.Spelling)
	return m.cfg.Theme.IncorrectColor.Render(str)
}

func (m quiz) RenderProgress() string {
	progressStr := "\nProgress: 0/0 (0%)\n\n"
	if m.attempts > 0 {
		progressStr = fmt.Sprintf("\nProgress: %d/%d (%.0f%%)\n",
			m.correct, m.attempts, float64(m.correct)/float64(m.attempts)*100)
	}
	return m.cfg.Theme.ProgressColor.Render(progressStr)
}

func (m quiz) RenderPrompt() string {
	promptStr, err := RenderTemplate(m.cfg.Prompt.Format, m.char)
	if err != nil {
		promptStr = "[template error]\n"
	}

	return m.cfg.Theme.PromptColor.Render(promptStr)
}

func (m quiz) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m quiz) RenderLeave() string {
	str := fmt.Sprintf("\nFinal score: %d/%d (%.0f%%)\nGoodbye!\n", m.correct, m.attempts, float64(m.correct)/float64(m.attempts)*100)
	return m.cfg.Theme.QuitColor.Render(str)
}

func (m quiz) RenderWelcome() string {
	str := fmt.Sprintf("Ayo v0.1 — Type '%s' to quit", m.cfg.Cmd.Exit)
	return m.cfg.Theme.WelcomeColor.Render(str)
}

func (m quiz) View() string {
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

func initialQuiz(cfg config.Config, chars []characters.Character) quiz {
	ti := textinput.New()
	ti.Placeholder = "Type here"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return quiz{
		cfg:       cfg,
		chars:     chars,
		char:      chars[rand.Intn(len(chars))],
		textInput: ti,
	}
}

func Quiz(cfg config.Config, chars []characters.Character) {
	p := tea.NewProgram(initialQuiz(cfg, chars))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
