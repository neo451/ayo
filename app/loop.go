package loop

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/neo451/alpha/internal/characters"
	"github.com/neo451/alpha/internal/config"
	"math/rand"
	"os"
	"strings"
	"text/template"
)

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

func Loop(cfg config.Config, characters []characters.Character) {
	scanner := bufio.NewScanner(os.Stdin)
	correct := 0
	attempts := 0

	fmt.Println("Language Character Quiz")
	fmt.Println("Type the spelling for each character. Type 'exit' to quit.")
	fmt.Printf("Total characters: %d\n\n", len(characters))

	// TODO: Ctrl-D
	for {
		// Select random character
		index := rand.Intn(len(characters))
		char := characters[index]
		prompt, err := RenderTemplate(cfg.Prompt.Format, char)

		if err != nil {
			panic("Wrong use of template")
		}
		fmt.Print(prompt)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		// Exit condition
		if strings.EqualFold(input, cfg.Cmd["exit"]) {
			break
		}

		attempts++
		if strings.EqualFold(input, char.Spelling) {
			fmt.Println(cfg.Prompt.Ok)
			correct++
		} else {
			fmt.Printf(cfg.Prompt.Err+"\n", char.Spelling)
		}
		ShowProgress(cfg, attempts, correct)
	}
}

func ShowProgress(cfg config.Config, attempts int, correct int) {
	// Show progress every 5 attempts
	if attempts%cfg.Progress.Frequency == 0 {
		fmt.Printf("\nProgress: %d/%d (%.0f%%)\n\n",
			correct, attempts, float64(correct)/float64(attempts)*100)
	}
}

// TODO: Final stats
func ShowFinal(cfg config.Config, attempts int, correct int) {
	if attempts > 0 {
		fmt.Printf("\nFinal score: %d/%d (%.0f%%)\n",
			correct, attempts, float64(correct)/float64(attempts)*100)
	}
	fmt.Println("Goodbye!")
}
